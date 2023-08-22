package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

type jsonFileCache string

func (cache jsonFileCache) load() (map[string]string, error) {
	data, err := os.ReadFile(string(cache))
	if err != nil {
		return nil, err
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

	val := values[key]
	if val == "" {
		return nil, autocert.ErrCacheMiss
	}
	data, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Put stores the data in the cache under the specified key.
// Underlying implementations may use any data storage format,
// as long as the reverse operation, Get, results in the original data.
func (cache jsonFileCache) Put(ctx context.Context, key string, data []byte) error {
	values, err := cache.load()
	if err != nil {
		return err
	}

	val := base64.StdEncoding.EncodeToString(data)
	values[key] = val

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
