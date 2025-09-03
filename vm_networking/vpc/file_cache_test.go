package networkvpc

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockedStorage struct {
	fileContent []byte
}

func NewMockedStorage() *MockedStorage {
	fileData := make([]byte, 1000)
	for i := 0; i < len(fileData); i++ {
		fileData[i] = byte(i % 256)
	}
	return &MockedStorage{
		fileContent: fileData,
	}
}

func (ms *MockedStorage) ReadFileChunk(path string, buffer []byte, index int64) (int, error) {
	if index == 112 {
		n := len(buffer)
		for i := 0; i < len(buffer); i++ {
			buffer[i] = ms.fileContent[index+int64(i)]
		}
		return n, nil
	}
	if index == 960 {
		n := 40
		for i := 0; i < 40; i++ {
			buffer[i] = ms.fileContent[index+int64(i)]
		}
		return n, io.ErrUnexpectedEOF
	}
	if index == 100 {
		n := len(buffer)
		for i := 0; i < len(buffer); i++ {
			buffer[i] = ms.fileContent[index+int64(i)]
		}
		return n, nil
	}
	return 0, errors.New("unexpected request in mocked test")
}

func Test_BufferAndIndexAreAtEndOfFile(t *testing.T) {
	cache := &ChunkCache{
		filePath: "",
		err:      nil,
		n:        1000,
		page:     7,
		buffer:   make([]byte, 100),
		storage:  NewMockedStorage(),
	}
	assert.False(t, cache.BufferAndIndexAreAtEndOfFile(712), "Expected not to be at end of file")

	cache = &ChunkCache{
		filePath: "",
		err:      io.EOF,
		n:        1000,
		page:     9,
		buffer:   make([]byte, 100),
		storage:  NewMockedStorage(),
	}
	assert.True(t, cache.BufferAndIndexAreAtEndOfFile(999), "Expected to be at end of file")
}

func Test_SlideBufferToIndex(t *testing.T) {
	cache := &ChunkCache{
		filePath: "",
		err:      nil,
		n:        1000,
		page:     -1,
		buffer:   make([]byte, 100),
		storage:  NewMockedStorage(),
	}
	cache.SlideBufferToIndex(112)
	buffer := cache.GetBuffered(115)
	for i := 0; i < cache.n; i++ {
		assert.Equal(t, buffer[i], byte((112+i)%256), "Expected the correct array value")
	}
	assert.Nil(t, cache.err, "No error expected")
	cache.SlideBufferToIndex(960)
	for i := 0; i < cache.n; i++ {
		assert.Equal(t, buffer[i], byte((960+i)%256), "Expected the correct array value")
	}
	assert.ErrorIs(t, cache.err, io.ErrUnexpectedEOF, "Expected EOF error")
}

func Test_getPage(t *testing.T) {
	cache := &ChunkCache{
		filePath: "",
		err:      nil,
		n:        1000,
		page:     -1,
		buffer:   make([]byte, 100),
		storage:  NewMockedStorage(),
	}
	assert.Equal(t, cache.getPage(150), int64(1), "Expected page 1 for index 150")
	assert.Equal(t, cache.getPage(1000), int64(10), "Expected page 10 for index 1000")
	assert.Equal(t, cache.getPage(100), int64(1), "Expected page 1 for index 100")
}

func Test_GetBuffered(t *testing.T) {
	cache := &ChunkCache{
		filePath: "",
		err:      nil,
		n:        1000,
		page:     -1,
		buffer:   make([]byte, 100),
		storage:  NewMockedStorage(),
	}
	buffer := cache.GetBuffered(115)
	for i := 0; i < cache.n; i++ {
		assert.Equal(t, buffer[i], byte((100+i)%256), "Expected the correct array value")
	}
}
