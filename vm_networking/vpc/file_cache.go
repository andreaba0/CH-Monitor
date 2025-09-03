package networkvpc

import (
	"errors"
	"io"
	"os"
)

type storageCache struct{}

func (s *storageCache) ReadFileChunk(path string, buffer []byte, index int64) (int, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	return fd.ReadAt(buffer, index)
}

type storageCacheService interface {
	ReadFileChunk(path string, buffer []byte, index int64) (int, error)
}

type ChunkCache struct {
	filePath string
	err      error
	n        int
	page     int64
	buffer   []byte
	storage  storageCacheService
}

func NewChunkCache(filePath string, bufferSize int) *ChunkCache {
	return &ChunkCache{
		filePath: filePath,
		err:      nil,
		n:        0,
		page:     -1,
		buffer:   make([]byte, bufferSize),
		storage:  new(storageCache),
	}
}

func (cache *ChunkCache) Error() error {
	return cache.err
}

func (cache *ChunkCache) getPage(index int64) int64 {
	return index / int64(len(cache.buffer))
}

// This method is used to align a page to
func (cache *ChunkCache) SlideBufferToIndex(index int64) {
	n, err := cache.storage.ReadFileChunk(cache.filePath, cache.buffer, index)
	cache.n = n
	cache.err = err
	cache.page = cache.getPage(index)
}

func (cache *ChunkCache) GetBuffered(index int64) []byte {
	if cache.page > -1 && cache.getPage(index) == cache.page {
		return cache.buffer
	}
	n, err := cache.storage.ReadFileChunk(cache.filePath, cache.buffer, index)
	cache.n = n
	cache.err = err
	cache.page = cache.getPage(index)
	return cache.buffer
}

func (cache *ChunkCache) BufferAndIndexAreAtEndOfFile(index int64) bool {
	if cache.getPage(index) < cache.page {
		return false
	}
	return index >= int64(cache.n)-1 && cache.EndFileReached()
}

func (cache *ChunkCache) EndFileReached() bool {
	return cache.err != nil && (errors.Is(cache.err, io.ErrUnexpectedEOF) || errors.Is(cache.err, io.EOF))
}
