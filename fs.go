package jsondb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewFS[T any](dir string) (Client[T], error) {
	// Check if dir exists and is a directory
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	return &fsClient[T]{dir: dir}, nil
}

type fsClient[T any] struct {
	dir string
}

func (f fsClient[T]) path(id string) string {
	return filepath.Join(f.dir, fmt.Sprintf("%s.json", id))
}

func (f *fsClient[T]) Write(id string, data T) error {
	byts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(f.path(id), byts, 0644)
}

func (f *fsClient[T]) Read(id string) (T, error) {
	var res T
	byts, err := os.ReadFile(f.path(id))
	if err != nil {
		if os.IsNotExist(err) {
			return res, NotFoundError{id: id}
		}
		return res, err
	}

	err = json.Unmarshal(byts, &res)
	return res, err
}

func (f *fsClient[T]) List() ([]T, error) {
	res := make([]T, 0)

	files, err := os.ReadDir(f.dir)
	if err != nil {
		return nil, err
	}

	// TODO: this would probably be faster using goroutines
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		id := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		obj, err := f.Read(id)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}
	return res, nil
}

func (f *fsClient[T]) Delete(id string) error {
	err := os.Remove(f.path(id))
	if err != nil {
		if os.IsNotExist(err) {
			return NotFoundError{id: id}
		}
		return err
	}
	return nil
}
