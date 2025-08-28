package networkvpc

import (
	"io"
	storage "vmm/storage"
)

type ChunkCache struct {
	filePath string
	err      error
	n        int
	page     int64
	buffer   []byte
}

func NewChunkCache(filePath string, bufferSize int) *ChunkCache {
	return &ChunkCache{
		filePath: filePath,
		err:      nil,
		n:        0,
		page:     0,
		buffer:   make([]byte, bufferSize),
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
	n, err := storage.ReadFileChunk(cache.filePath, cache.buffer, index)
	cache.n = n
	cache.err = err
}

func (cache *ChunkCache) GetBuffered(index int64) []byte {
	if cache.getPage(index) == cache.page {
		return cache.buffer
	}
	n, err := storage.ReadFileChunk(cache.filePath, cache.buffer, index)
	cache.n = n
	cache.err = err
	cache.page = cache.getPage(index)
	return cache.buffer
}

func (cache *ChunkCache) BufferAndIndexAreAtEndOfFile(index int64) bool {
	if cache.getPage(index) < cache.page {
		return false
	}
	return (index%int64(len(cache.buffer))) >= int64(cache.n) && cache.EndFileReached()
}

func (cache *ChunkCache) EndFileReached() bool {
	return cache.err != nil && (cache.err == io.ErrUnexpectedEOF || cache.err == io.EOF)
}
