package main

import "time"

type TransferType string

const (
	TransferUpload   TransferType = "upload"
	TransferDownload TransferType = "download"
)

type TransferStatus string

const (
	StatusPending    TransferStatus = "pending"
	StatusInProgress TransferStatus = "in_progress"
	StatusCompleted  TransferStatus = "completed"
	StatusFailed     TransferStatus = "failed"
	StatusCancelled  TransferStatus = "cancelled"
)

type Connection struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	UseSSL    bool   `json:"useSSL"`
	PathStyle bool   `json:"pathStyle"`
}

type S3Object struct {
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	IsFolder     bool      `json:"isFolder"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	StorageClass string    `json:"storageClass"`
}

type Settings struct {
	MaxParallel int `json:"maxParallel"`
}

type Transfer struct {
	ID         string         `json:"id"`
	FileName   string         `json:"fileName"`
	Type       TransferType   `json:"type"`
	Status     TransferStatus `json:"status"`
	BytesTotal int64          `json:"bytesTotal"`
	BytesDone  int64          `json:"bytesDone"`
	Percentage float64        `json:"percentage"`
	Error      string         `json:"error"`
}
