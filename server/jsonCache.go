package server

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type jsonFileCache string

func (cache jsonFileCache) load() (map[string]string, error) {
	data, err := os.ReadFile(string(cache))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			data = []byte("{}")
		} else {
			return nil, err
		}
	}

	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}

	return values, nil
}

func (cache jsonFileCache) save(values map[string]string) error {
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(string(cache), data, 0600); err != nil {
		return err
	}

	return nil
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (cache jsonFileCache) Get(ctx context.Context, key string) ([]byte, error) {
	values, err := cache.load()
	if err != nil {
		return nil, err
	}

	data := values[key]
	if data == "" {
		return nil, autocert.ErrCacheMiss
	}

	return []byte(data), nil
}

// Put stores the data in the cache under the specified key.
// Underlying implementations may use any data storage format,
// as long as the reverse operation, Get, results in the original data.
func (cache jsonFileCache) Put(ctx context.Context, key string, data []byte) error {
	values, err := cache.load()
	if err != nil {
		return err
	}

	values[key] = string(data)

	return cache.save(values)
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (cache jsonFileCache) Delete(ctx context.Context, key string) error {
	values, err := cache.load()
	if err != nil {
		return err
	}

	delete(values, key)

	return cache.save(values)
}

var _ autocert.Cache = jsonFileCache("")
