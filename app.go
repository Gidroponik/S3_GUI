package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	config     *ConfigStore
	s3client   *S3ClientWrapper
	transfers  *TransferManager
	connecting bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.config = NewConfigStore()
	settings := a.config.LoadSettings()
	a.transfers = NewTransferManager(ctx, settings.MaxParallel)
}

func (a *App) shutdown(ctx context.Context) {
	a.Disconnect()
}

// --- Connection management ---

func (a *App) GetConnections() ([]Connection, error) {
	return a.config.LoadAll()
}

func (a *App) SaveConnection(conn Connection) error {
	return a.config.Save(conn)
}

func (a *App) DeleteConnection(id string) error {
	return a.config.Delete(id)
}

// --- S3 operations ---

func (a *App) Connect(id string) error {
	if a.connecting {
		return fmt.Errorf("connection already in progress")
	}
	a.connecting = true
	defer func() { a.connecting = false }()

	log.Printf("[Connect] connecting to id=%s", id)
	conns, err := a.config.LoadAll()
	if err != nil {
		log.Printf("[Connect] LoadAll error: %v", err)
		return err
	}
	log.Printf("[Connect] loaded %d connections", len(conns))
	for _, c := range conns {
		log.Printf("[Connect] checking id=%s name=%s", c.ID, c.Name)
		if c.ID == id {
			log.Printf("[Connect] found match, host=%s:%d bucket=%s ssl=%v pathStyle=%v", c.Host, c.Port, c.Bucket, c.UseSSL, c.PathStyle)
			wrapper, err := NewS3ClientWrapper(c)
			if err != nil {
				log.Printf("[Connect] NewS3ClientWrapper error: %v", err)
				return err
			}
			if err := wrapper.TestConnection(a.ctx); err != nil {
				log.Printf("[Connect] TestConnection error: %v", err)
				return err
			}
			log.Printf("[Connect] success!")
			a.s3client = wrapper
			return nil
		}
	}
	return fmt.Errorf("connection not found: %s", id)
}

func (a *App) Disconnect() {
	a.s3client = nil
}

func (a *App) TestConnection(conn Connection) error {
	log.Printf("[TestConnection] host=%s:%d bucket=%s ssl=%v pathStyle=%v", conn.Host, conn.Port, conn.Bucket, conn.UseSSL, conn.PathStyle)
	wrapper, err := NewS3ClientWrapper(conn)
	if err != nil {
		log.Printf("[TestConnection] NewS3ClientWrapper error: %v", err)
		return err
	}
	err = wrapper.TestConnection(a.ctx)
	if err != nil {
		log.Printf("[TestConnection] error: %v", err)
	} else {
		log.Printf("[TestConnection] success!")
	}
	return err
}

func (a *App) ListObjects(prefix string) ([]S3Object, error) {
	if a.s3client == nil {
		return nil, fmt.Errorf("not connected")
	}
	return a.s3client.ListObjects(a.ctx, prefix)
}

func (a *App) DeleteObjects(keys []string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	return a.s3client.DeleteObjects(a.ctx, keys)
}

func (a *App) CreateFolder(prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	return a.s3client.CreateFolder(a.ctx, prefix)
}

func (a *App) GetPresignedURL(key string) (string, error) {
	if a.s3client == nil {
		return "", fmt.Errorf("not connected")
	}
	return a.s3client.GetPresignedURL(a.ctx, key, 1*time.Hour)
}

func (a *App) GetDirectURL(key string) (string, error) {
	if a.s3client == nil {
		return "", fmt.Errorf("not connected")
	}
	return a.s3client.GetDirectURL(key), nil
}

func (a *App) DeletePrefix(prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	return a.s3client.DeletePrefix(a.ctx, prefix)
}

func (a *App) UploadDroppedFiles(paths []string, prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			log.Printf("[UploadDroppedFiles] stat error for %s: %v", p, err)
			continue
		}
		if info.IsDir() {
			if err := a.transfers.QueueFolderUpload(a.s3client, p, prefix); err != nil {
				log.Printf("[UploadDroppedFiles] folder upload error: %v", err)
			}
		} else {
			a.transfers.QueueUpload(a.s3client, p, prefix)
		}
	}
	return nil
}

// --- Transfers ---

func (a *App) UploadFiles(prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select files to upload",
	})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	for _, f := range files {
		a.transfers.QueueUpload(a.s3client, f, prefix)
	}
	return nil
}

func (a *App) DownloadFiles(keys []string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select download folder",
	})
	if err != nil {
		return err
	}
	if dir == "" {
		return nil
	}
	for _, key := range keys {
		a.transfers.QueueDownload(a.s3client, key, dir)
	}
	return nil
}

func (a *App) UploadFolder(prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select folder to upload",
	})
	if err != nil {
		return err
	}
	if dir == "" {
		return nil
	}
	return a.transfers.QueueFolderUpload(a.s3client, dir, prefix)
}

func (a *App) DownloadPrefix(prefix string) error {
	if a.s3client == nil {
		return fmt.Errorf("not connected")
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select download folder",
	})
	if err != nil {
		return err
	}
	if dir == "" {
		return nil
	}
	return a.transfers.QueuePrefixDownload(a.s3client, prefix, dir)
}

func (a *App) CancelTransfer(id string) {
	a.transfers.Cancel(id)
}

func (a *App) ClearCompletedTransfers() {
	a.transfers.ClearCompleted()
}

func (a *App) GetTransfers() []Transfer {
	return a.transfers.GetAll()
}

// --- Settings ---

func (a *App) GetSettings() Settings {
	return a.config.LoadSettings()
}

func (a *App) SaveSettings(s Settings) error {
	if err := a.config.SaveSettings(s); err != nil {
		return err
	}
	a.transfers.SetConcurrency(s.MaxParallel)
	return nil
}

// --- Export / Import Connections ---

func (a *App) ExportConnections() error {
	conns, err := a.config.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load connections: %w", err)
	}
	if len(conns) == 0 {
		return fmt.Errorf("no connections to export")
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export Connections",
		DefaultFilename: "s3-connections.s3b",
		Filters: []runtime.FileFilter{
			{DisplayName: "S3 Backup Files", Pattern: "*.s3b"},
		},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	data, err := json.Marshal(conns)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return os.WriteFile(path, []byte(encoded), 0600)
}

func (a *App) ImportConnections() (int, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Import Connections",
		Filters: []runtime.FileFilter{
			{DisplayName: "S3 Backup Files", Pattern: "*.s3b"},
		},
	})
	if err != nil {
		return 0, err
	}
	if path == "" {
		return 0, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}
	data, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return 0, fmt.Errorf("invalid backup file: %w", err)
	}
	var imported []Connection
	if err := json.Unmarshal(data, &imported); err != nil {
		return 0, fmt.Errorf("corrupted backup data: %w", err)
	}
	existing, err := a.config.LoadAll()
	if err != nil {
		existing = []Connection{}
	}
	isDuplicate := func(c Connection) bool {
		for _, e := range existing {
			if e.Name == c.Name && e.Host == c.Host && e.Bucket == c.Bucket {
				return true
			}
		}
		return false
	}
	added := 0
	for _, c := range imported {
		if isDuplicate(c) {
			continue
		}
		c.ID = uuid.New().String()
		if err := a.config.Save(c); err != nil {
			log.Printf("[ImportConnections] save error: %v", err)
			continue
		}
		existing = append(existing, c)
		added++
	}
	return added, nil
}
