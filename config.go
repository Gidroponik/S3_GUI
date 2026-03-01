package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type ConfigStore struct {
	path         string
	settingsPath string
}

func NewConfigStore() *ConfigStore {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	dir := filepath.Join(appData, "S3BucketGUI")
	os.MkdirAll(dir, 0700)
	return &ConfigStore{
		path:         filepath.Join(dir, "connections.dat"),
		settingsPath: filepath.Join(dir, "settings.json"),
	}
}

func (c *ConfigStore) LoadAll() ([]Connection, error) {
	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Connection{}, nil
		}
		return nil, err
	}
	plaintext, err := decrypt(string(data))
	if err != nil {
		return nil, err
	}
	var conns []Connection
	if err := json.Unmarshal(plaintext, &conns); err != nil {
		return nil, err
	}
	return conns, nil
}

func (c *ConfigStore) saveAll(conns []Connection) error {
	data, err := json.Marshal(conns)
	if err != nil {
		return err
	}
	encrypted, err := encrypt(data)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, []byte(encrypted), 0600)
}

func (c *ConfigStore) Save(conn Connection) error {
	conns, err := c.LoadAll()
	if err != nil {
		conns = []Connection{}
	}
	if conn.ID == "" {
		conn.ID = uuid.New().String()
	}
	found := false
	for i, existing := range conns {
		if existing.ID == conn.ID {
			conns[i] = conn
			found = true
			break
		}
	}
	if !found {
		conns = append(conns, conn)
	}
	return c.saveAll(conns)
}

func (c *ConfigStore) Delete(id string) error {
	conns, err := c.LoadAll()
	if err != nil {
		return err
	}
	filtered := make([]Connection, 0, len(conns))
	for _, conn := range conns {
		if conn.ID != id {
			filtered = append(filtered, conn)
		}
	}
	return c.saveAll(filtered)
}

func (c *ConfigStore) LoadSettings() Settings {
	defaults := Settings{MaxParallel: 3}
	data, err := os.ReadFile(c.settingsPath)
	if err != nil {
		return defaults
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return defaults
	}
	if s.MaxParallel < 1 {
		s.MaxParallel = 1
	}
	if s.MaxParallel > 10 {
		s.MaxParallel = 10
	}
	return s
}

func (c *ConfigStore) SaveSettings(s Settings) error {
	if s.MaxParallel < 1 {
		s.MaxParallel = 1
	}
	if s.MaxParallel > 10 {
		s.MaxParallel = 10
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.settingsPath, data, 0600)
}
