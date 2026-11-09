package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Storage handles saving and loading student data to/from disk.
type Storage struct {
	filePath string
}

// NewStorage initializes a new Storage instance with the given file path.
func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

// Save writes the student records to the specified file in JSON format.
func (s *Storage) Save(students map[int]Student) error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

// Load reads the student records from the specified file.
// If the file does not exist or is empty, it returns an empty map without error.
func (s *Storage) Load() (map[int]Student, error) {
	students := make(map[int]Student)

	if _, err := os.Stat(s.filePath); errors.Is(err, os.ErrNotExist) {
		return students, nil
	}

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return students, nil
	}

	if err := json.Unmarshal(data, &students); err != nil {
		return nil, err
	}

	return students, nil
}
