package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TransferManager struct {
	ctx       context.Context
	mu        sync.Mutex
	transfers map[string]*transferEntry
	sem       chan struct{}
}

type transferEntry struct {
	Transfer
	cancel context.CancelFunc
}

func NewTransferManager(ctx context.Context, maxParallel int) *TransferManager {
	if maxParallel < 1 {
		maxParallel = 3
	}
	return &TransferManager{
		ctx:       ctx,
		transfers: make(map[string]*transferEntry),
		sem:       make(chan struct{}, maxParallel),
	}
}

func (tm *TransferManager) SetConcurrency(n int) {
	if n < 1 {
		n = 1
	}
	if n > 10 {
		n = 10
	}
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.sem = make(chan struct{}, n)
}

func (tm *TransferManager) GetAll() []Transfer {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	result := make([]Transfer, 0, len(tm.transfers))
	for _, e := range tm.transfers {
		result = append(result, e.Transfer)
	}
	return result
}

func (tm *TransferManager) Cancel(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if e, ok := tm.transfers[id]; ok {
		if e.cancel != nil {
			e.cancel()
		}
		e.Status = StatusCancelled
		tm.emitUpdate(e.Transfer)
	}
}

func (tm *TransferManager) ClearCompleted() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for id, e := range tm.transfers {
		if e.Status == StatusCompleted || e.Status == StatusFailed || e.Status == StatusCancelled {
			delete(tm.transfers, id)
		}
	}
	runtime.EventsEmit(tm.ctx, "transfer:cleared")
}

func (tm *TransferManager) emitUpdate(t Transfer) {
	runtime.EventsEmit(tm.ctx, "transfer:progress", t)
}

func (tm *TransferManager) QueueUpload(client *S3ClientWrapper, localPath, prefix string) {
	id := uuid.New().String()
	fileName := filepath.Base(localPath)
	key := prefix + fileName

	entry := &transferEntry{
		Transfer: Transfer{
			ID:       id,
			FileName: fileName,
			Type:     TransferUpload,
			Status:   StatusPending,
		},
	}

	tm.mu.Lock()
	tm.transfers[id] = entry
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)

	go tm.runUpload(client, entry, localPath, key)
}

func (tm *TransferManager) runUpload(client *S3ClientWrapper, entry *transferEntry, localPath, key string) {
	tm.sem <- struct{}{}
	defer func() { <-tm.sem }()

	ctx, cancel := context.WithCancel(tm.ctx)
	tm.mu.Lock()
	entry.cancel = cancel
	entry.Status = StatusInProgress
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)
	defer cancel()

	f, err := os.Open(localPath)
	if err != nil {
		tm.failTransfer(entry, err)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		tm.failTransfer(entry, err)
		return
	}

	tm.mu.Lock()
	entry.BytesTotal = info.Size()
	tm.mu.Unlock()

	pr := &progressReader{
		reader:  f,
		total:   info.Size(),
		onProgress: tm.throttledProgress(entry),
	}

	uploader := manager.NewUploader(client.client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(client.conn.Bucket),
		Key:    aws.String(key),
		Body:   pr,
	})

	if err != nil {
		if ctx.Err() != nil {
			return
		}
		tm.failTransfer(entry, err)
		return
	}

	tm.mu.Lock()
	entry.Status = StatusCompleted
	entry.BytesDone = entry.BytesTotal
	entry.Percentage = 100
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)
	runtime.EventsEmit(tm.ctx, "transfer:complete", entry.Transfer)
}

func (tm *TransferManager) QueueDownload(client *S3ClientWrapper, key, localDir string) {
	id := uuid.New().String()
	fileName := filepath.Base(key)

	entry := &transferEntry{
		Transfer: Transfer{
			ID:       id,
			FileName: fileName,
			Type:     TransferDownload,
			Status:   StatusPending,
		},
	}

	tm.mu.Lock()
	tm.transfers[id] = entry
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)

	go tm.runDownload(client, entry, key, filepath.Join(localDir, fileName))
}

func (tm *TransferManager) runDownload(client *S3ClientWrapper, entry *transferEntry, key, localPath string) {
	tm.sem <- struct{}{}
	defer func() { <-tm.sem }()

	ctx, cancel := context.WithCancel(tm.ctx)
	tm.mu.Lock()
	entry.cancel = cancel
	entry.Status = StatusInProgress
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)
	defer cancel()

	head, err := client.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(client.conn.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		tm.failTransfer(entry, err)
		return
	}

	tm.mu.Lock()
	entry.BytesTotal = aws.ToInt64(head.ContentLength)
	tm.mu.Unlock()

	os.MkdirAll(filepath.Dir(localPath), 0755)
	f, err := os.Create(localPath)
	if err != nil {
		tm.failTransfer(entry, err)
		return
	}
	defer f.Close()

	pw := &progressWriterAt{
		writer: f,
		total:  aws.ToInt64(head.ContentLength),
		onProgress: tm.throttledProgress(entry),
	}

	downloader := manager.NewDownloader(client.client)
	_, err = downloader.Download(ctx, pw, &s3.GetObjectInput{
		Bucket: aws.String(client.conn.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if ctx.Err() != nil {
			return
		}
		tm.failTransfer(entry, err)
		return
	}

	tm.mu.Lock()
	entry.Status = StatusCompleted
	entry.BytesDone = entry.BytesTotal
	entry.Percentage = 100
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)
	runtime.EventsEmit(tm.ctx, "transfer:complete", entry.Transfer)
}

func (tm *TransferManager) QueueFolderUpload(client *S3ClientWrapper, localDir, prefix string) error {
	baseName := filepath.Base(localDir)
	targetPrefix := prefix + baseName + "/"
	return filepath.WalkDir(localDir, func(fpath string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(localDir, fpath)
		relKey := strings.ReplaceAll(rel, "\\", "/")
		tm.queueUploadWithKey(client, fpath, targetPrefix+relKey)
		return nil
	})
}

func (tm *TransferManager) queueUploadWithKey(client *S3ClientWrapper, localPath, key string) {
	id := uuid.New().String()
	fileName := filepath.Base(localPath)

	entry := &transferEntry{
		Transfer: Transfer{
			ID:       id,
			FileName: fileName,
			Type:     TransferUpload,
			Status:   StatusPending,
		},
	}

	tm.mu.Lock()
	tm.transfers[id] = entry
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)

	go tm.runUpload(client, entry, localPath, key)
}

func (tm *TransferManager) QueuePrefixDownload(client *S3ClientWrapper, prefix, localDir string) error {
	objects, err := client.ListAllObjects(tm.ctx, prefix)
	if err != nil {
		return err
	}
	for _, obj := range objects {
		if obj.IsFolder || obj.Size == 0 {
			continue
		}
		rel := strings.TrimPrefix(obj.Key, prefix)
		destPath := filepath.Join(localDir, filepath.FromSlash(rel))
		os.MkdirAll(filepath.Dir(destPath), 0755)
		tm.QueueDownload(client, obj.Key, filepath.Dir(destPath))
	}
	return nil
}

func (tm *TransferManager) failTransfer(entry *transferEntry, err error) {
	tm.mu.Lock()
	entry.Status = StatusFailed
	entry.Error = err.Error()
	tm.mu.Unlock()
	tm.emitUpdate(entry.Transfer)
}

func (tm *TransferManager) throttledProgress(entry *transferEntry) func(done int64) {
	var lastEmit time.Time
	return func(done int64) {
		tm.mu.Lock()
		entry.BytesDone = done
		if entry.BytesTotal > 0 {
			entry.Percentage = float64(done) / float64(entry.BytesTotal) * 100
		}
		t := entry.Transfer
		tm.mu.Unlock()

		now := time.Now()
		if now.Sub(lastEmit) >= 100*time.Millisecond {
			lastEmit = now
			tm.emitUpdate(t)
		}
	}
}

// progressReader wraps an io.Reader to track upload progress
type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	onProgress func(done int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.read += int64(n)
	if pr.onProgress != nil {
		pr.onProgress(pr.read)
	}
	return n, err
}

// progressWriterAt wraps an io.WriterAt to track download progress
type progressWriterAt struct {
	writer     io.WriterAt
	total      int64
	mu         sync.Mutex
	written    int64
	onProgress func(done int64)
}

func (pw *progressWriterAt) WriteAt(p []byte, off int64) (int, error) {
	n, err := pw.writer.WriteAt(p, off)
	pw.mu.Lock()
	pw.written += int64(n)
	done := pw.written
	pw.mu.Unlock()
	if pw.onProgress != nil {
		pw.onProgress(done)
	}
	return n, err
}

// Ensure progressReader satisfies io.Reader
var _ io.Reader = (*progressReader)(nil)

// Ensure progressWriterAt satisfies io.WriterAt
var _ io.WriterAt = (*progressWriterAt)(nil)

// We need a Seek method for the upload manager
func (pr *progressReader) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := pr.reader.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, fmt.Errorf("seek not supported")
}
