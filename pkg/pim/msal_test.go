/*
Copyright © 2026 netr0m <netr0m@pm.me>
*/
package pim

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
	"github.com/stretchr/testify/assert"
)

type fakeUnmarshaler struct {
	called bool
	got    []byte
	err    error
}

func (f *fakeUnmarshaler) Unmarshal(data []byte) error {
	f.called = true
	f.got = data
	return f.err
}

type fakeMarshaler struct {
	data []byte
	err  error
}

func (f *fakeMarshaler) Marshal() ([]byte, error) {
	return f.data, f.err
}

func TestFileCacheReplaceNoFile(t *testing.T) {
	fc := &fileCache{path: filepath.Join(t.TempDir(), "does-not-exist.json")}
	u := &fakeUnmarshaler{}

	err := fc.Replace(context.Background(), u, cache.ReplaceHints{})

	assert.NoError(t, err)
	assert.False(t, u.called, "expected Unmarshal not to be called when the cache file doesn't exist")
}

func TestFileCacheReplaceUnmarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	if err := os.WriteFile(path, []byte("not valid"), 0o600); err != nil {
		t.Fatalf("failed to seed cache file: %v", err)
	}
	fc := &fileCache{path: path}
	u := &fakeUnmarshaler{err: errors.New("boom")}

	err := fc.Replace(context.Background(), u, cache.ReplaceHints{})

	assert.Error(t, err)
	assert.True(t, u.called)
}

func TestFileCacheExportThenReplace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "cache.json")
	fc := &fileCache{path: path}
	want := []byte(`{"token":"dummy"}`)

	err := fc.Export(context.Background(), &fakeMarshaler{data: want}, cache.ExportHints{})
	assert.NoError(t, err)

	onDisk, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, want, onDisk)

	u := &fakeUnmarshaler{}
	err = fc.Replace(context.Background(), u, cache.ReplaceHints{})
	assert.NoError(t, err)
	assert.True(t, u.called)
	assert.Equal(t, want, u.got)
}

func TestFileCacheExportMarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cache.json")
	fc := &fileCache{path: path}

	err := fc.Export(context.Background(), &fakeMarshaler{err: errors.New("boom")}, cache.ExportHints{})

	assert.Error(t, err)
	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr), "expected no cache file to be written on marshal error")
}

func TestClearGraphTokenCacheRemovesExistingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CACHE_HOME", dir)

	path, err := graphTokenCachePath()
	assert.NoError(t, err)
	assert.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	assert.NoError(t, os.WriteFile(path, []byte(`{"token":"dummy"}`), 0o600))

	assert.NoError(t, ClearGraphTokenCache())

	_, statErr := os.Stat(path)
	assert.True(t, os.IsNotExist(statErr), "expected token cache file to be removed")
}

func TestClearGraphTokenCacheNoFileIsNotAnError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CACHE_HOME", dir)

	assert.NoError(t, ClearGraphTokenCache())
}
