package utils

import (
	"encoding/gob"
	"os"
)

func ReadGobFile[T any](path string) (T, error) {
	var zero T

	fd, err := os.Open(path)
	if err != nil {
		return zero, err
	}
	defer fd.Close()

	decoder := gob.NewDecoder(fd)

	var result T
	err = decoder.Decode(&result)
	return result, err
}

func WriteGobFile[T any](path string, db T) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	encoder := gob.NewEncoder(fd)
	return encoder.Encode(db)
}

func AppendOrCreateToFile(path string, row []byte) (int, error) {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	return fd.Write(row)
}

func CreateFile(path string) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func ReadChunkFromFile(path string, buffer []byte, index int64) (int, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	return fd.ReadAt(buffer, index)
}
