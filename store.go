package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	_ Store = (*InMemoryStore)(nil)
	_ Store = (*FilesStore)(nil)
)

type Store interface {
	// Get retrieves the URL for a given alias.
	Get(alias string) (string, bool)
	// Set stores a URL with the given alias.
	Set(alias, url string) error
	// All lists all stored aliases and their URLs.
	All() map[string]string
}

type InMemoryStore struct {
	urlStore map[string]string
	mu       sync.Mutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		urlStore: make(map[string]string),
	}
}

func (s *InMemoryStore) Get(alias string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, exists := s.urlStore[alias]
	return url, exists
}

func (s *InMemoryStore) Set(alias, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urlStore[alias] = url
	return nil
}

func (s *InMemoryStore) All() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy of the map to avoid concurrent access issues
	copy := make(map[string]string, len(s.urlStore))
	for k, v := range s.urlStore {
		copy[k] = v
	}
	return copy
}

type FilesStore struct {
	directoryPath string
}

func NewFilesStore(directoryPath string) *FilesStore {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	storeDir := filepath.Join(userConfigDir, ApplicationName)
	if _, err := os.Stat(storeDir); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(storeDir, 0755)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	return &FilesStore{
		directoryPath: filepath.Join(userConfigDir, ApplicationName),
	}
}

func (s *FilesStore) Get(alias string) (string, bool) {
	filePath := filepath.Join(s.directoryPath, alias)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (s *FilesStore) Set(alias, url string) error {
	if strings.ContainsRune(alias, filepath.Separator) {
		return fmt.Errorf("invalid alias: %s", alias)
	}
	if strings.ContainsRune(alias, '.') {
		return fmt.Errorf("invalid alias: %s", alias)
	}
	filePath := filepath.Clean(filepath.Join(s.directoryPath, alias))
	return os.WriteFile(filePath, []byte(url), 0644)
}

func (s *FilesStore) All() map[string]string {
	files, err := os.ReadDir(s.directoryPath)
	if err != nil {
		panic(err)
	}

	urls := make(map[string]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.directoryPath, file.Name()))
		if err != nil {
			continue
		}
		urls[file.Name()] = string(data)
	}
	return urls
}
