package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/denisbrodbeck/machineid"
	"golang.org/x/crypto/pbkdf2"
)

var cachedKey []byte

func deriveKey() ([]byte, error) {
	if cachedKey != nil {
		return cachedKey, nil
	}
	mid, err := machineid.ProtectedID("S3BucketGUI")
	if err != nil {
		return nil, fmt.Errorf("failed to get machine id: %w", err)
	}
	salt := sha256.Sum256([]byte("S3BucketGUI-salt-" + mid))
	cachedKey = pbkdf2.Key([]byte(mid), salt[:], 100_000, 32, sha256.New)
	return cachedKey, nil
}

func encrypt(plaintext []byte) (string, error) {
	key, err := deriveKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encoded string) ([]byte, error) {
	key, err := deriveKey()
	if err != nil {
		return nil, err
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
