package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

type S3ClientWrapper struct {
	client *s3.Client
	conn   Connection
}

// stripSDKHeaders removes non-essential SDK headers before signing.
// When S3 sits behind a reverse proxy (Cloudflare, Caddy, nginx, etc.),
// proxies may modify or strip headers like Accept-Encoding, Amz-Sdk-*
// that the AWS SDK includes in the signed headers. If the server receives
// different header values than what was signed, signature verification fails.
// This middleware strips those headers so they are never signed.
func stripSDKHeaders() func(*middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		err := stack.Finalize.Insert(
			middleware.FinalizeMiddlewareFunc("StripSDKHeaders",
				func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
					if req, ok := input.Request.(*smithyhttp.Request); ok {
						req.Header.Del("Accept-Encoding")
						req.Header.Del("Amz-Sdk-Invocation-Id")
						req.Header.Del("Amz-Sdk-Request")
					}
					return handler.HandleFinalize(ctx, input)
				},
			),
			"Signing", middleware.Before,
		)
		if err != nil {
			// "Signing" step not found (e.g. presign flow) — skip silently
			return nil
		}
		return nil
	}
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func NewS3ClientWrapper(conn Connection) (*S3ClientWrapper, error) {
	scheme := "http"
	if conn.UseSSL {
		scheme = "https"
	}
	port := conn.Port
	if port == 0 {
		if conn.UseSSL {
			port = 443
		} else {
			port = 80
		}
	}

	var endpoint string
	if (conn.UseSSL && port == 443) || (!conn.UseSSL && port == 80) {
		endpoint = fmt.Sprintf("%s://%s", scheme, conn.Host)
	} else {
		endpoint = fmt.Sprintf("%s://%s:%d", scheme, conn.Host, port)
	}

	region := conn.Region
	if region == "" {
		region = "us-east-1"
	}

	log.Printf("[S3Client] endpoint=%s region=%s pathStyle=%v accessKey=%s secretKeyLen=%d",
		endpoint, region, conn.PathStyle, maskKey(conn.AccessKey), len(conn.SecretKey))

	httpClient := awshttp.NewBuildableClient().WithTransportOptions(func(t *http.Transport) {
		t.DialContext = (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
		t.TLSHandshakeTimeout = 10 * time.Second
		t.ResponseHeaderTimeout = 15 * time.Second
		t.ExpectContinueTimeout = 2 * time.Second
		t.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	})

	client := s3.New(s3.Options{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			conn.AccessKey, conn.SecretKey, "",
		),
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: conn.PathStyle,
		HTTPClient:   httpClient,
		HTTPSignerV4: v4.NewSigner(func(o *v4.SignerOptions) {
			o.DisableURIPathEscaping = true
		}),
		APIOptions: []func(*middleware.Stack) error{
			stripSDKHeaders(),
		},
	})

	return &S3ClientWrapper{client: client, conn: conn}, nil
}

func (w *S3ClientWrapper) TestConnection(ctx context.Context) error {
	testCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Try ListObjectsV2 with MaxKeys=1 — more reliable than HeadBucket
	// on custom S3 implementations (MinIO, Wasabi, etc.)
	_, err := w.client.ListObjectsV2(testCtx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(w.conn.Bucket),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	return nil
}

func (w *S3ClientWrapper) ListObjects(ctx context.Context, prefix string) ([]S3Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(w.conn.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	var objects []S3Object
	paginator := s3.NewListObjectsV2Paginator(w.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, cp := range page.CommonPrefixes {
			folderKey := aws.ToString(cp.Prefix)
			name := strings.TrimSuffix(strings.TrimPrefix(folderKey, prefix), "/")
			objects = append(objects, S3Object{
				Key:      folderKey,
				Name:     name,
				IsFolder: true,
			})
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			if key == prefix {
				continue
			}
			name := path.Base(key)
			objects = append(objects, S3Object{
				Key:          key,
				Name:         name,
				IsFolder:     false,
				Size:         aws.ToInt64(obj.Size),
				LastModified: aws.ToTime(obj.LastModified),
				StorageClass: string(obj.StorageClass),
			})
		}
	}

	return objects, nil
}

func (w *S3ClientWrapper) DeleteObjects(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// S3 API limits DeleteObjects to 1000 keys per request
	const batchSize = 1000
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]

		objectIds := make([]types.ObjectIdentifier, len(batch))
		for j, key := range batch {
			objectIds[j] = types.ObjectIdentifier{
				Key: aws.String(key),
			}
		}

		_, err := w.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(w.conn.Bucket),
			Delete: &types.Delete{
				Objects: objectIds,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *S3ClientWrapper) DeletePrefix(ctx context.Context, prefix string) error {
	objects, err := w.ListAllObjects(ctx, prefix)
	if err != nil {
		return fmt.Errorf("failed to list objects for deletion: %w", err)
	}

	// Collect keys, always include the folder marker itself
	keySet := make(map[string]struct{})
	for _, obj := range objects {
		keySet[obj.Key] = struct{}{}
	}
	keySet[prefix] = struct{}{}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}

	return w.DeleteObjects(ctx, keys)
}

func (w *S3ClientWrapper) CreateFolder(ctx context.Context, prefix string) error {
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	_, err := w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(w.conn.Bucket),
		Key:    aws.String(prefix),
		Body:   strings.NewReader(""),
	})
	return err
}

func (w *S3ClientWrapper) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(w.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(w.conn.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return req.URL, nil
}

func (w *S3ClientWrapper) GetDirectURL(key string) string {
	scheme := "http"
	if w.conn.UseSSL {
		scheme = "https"
	}
	port := w.conn.Port
	if port == 0 {
		if w.conn.UseSSL {
			port = 443
		} else {
			port = 80
		}
	}
	defaultPort := 80
	if w.conn.UseSSL {
		defaultPort = 443
	}
	portPart := ""
	if port != defaultPort {
		portPart = fmt.Sprintf(":%d", port)
	}
	if w.conn.PathStyle {
		return fmt.Sprintf("%s://%s%s/%s/%s", scheme, w.conn.Host, portPart, w.conn.Bucket, key)
	}
	return fmt.Sprintf("%s://%s.%s%s/%s", scheme, w.conn.Bucket, w.conn.Host, portPart, key)
}

func (w *S3ClientWrapper) ListAllObjects(ctx context.Context, prefix string) ([]S3Object, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(w.conn.Bucket),
		Prefix: aws.String(prefix),
	}

	var objects []S3Object
	paginator := s3.NewListObjectsV2Paginator(w.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			objects = append(objects, S3Object{
				Key:          key,
				Name:         path.Base(key),
				IsFolder:     strings.HasSuffix(key, "/"),
				Size:         aws.ToInt64(obj.Size),
				LastModified: aws.ToTime(obj.LastModified),
				StorageClass: string(obj.StorageClass),
			})
		}
	}

	return objects, nil
}
