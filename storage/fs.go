package storage

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func CreateFile(path string) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func CreateFileIfNotExists(filePath string) error {
	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func WriteFileChunk(path string, byteIndex int64, chunk io.Reader) error {
	fd, err := os.OpenFile(path, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Seek(byteIndex, 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(fd, chunk)
	if err != nil {
		return err
	}
	return nil
}

func ReadFileChunk(path string, buffer []byte, offset int64) (int, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	return fd.ReadAt(buffer, offset)
}

func RenameFile(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func ReadYaml[T any](path string) (T, error) {
	var res T
	dataByte, err := os.ReadFile(path)
	if err != nil {
		return res, err
	}
	err = yaml.Unmarshal(dataByte, &res)
	if err != nil {
		return res, err
	}
	return res, err
}

func ReadJson[T any](path string) (T, error) {
	var res T
	dataByte, err := os.ReadFile(path)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(dataByte, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func WriteJson[T any](path string, content T) error {
	dataByte, err := json.Marshal(content)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(dataByte)
	return err
}

func CreateFolderRec(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

type FolderEntry struct {
	Name     string
	IsFolder bool
}

func ListFolder(path string) ([]FolderEntry, error) {
	var res []FolderEntry = []FolderEntry{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return []FolderEntry{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			res = append(res, FolderEntry{
				IsFolder: true,
				Name:     entry.Name(),
			})
		} else {
			res = append(res, FolderEntry{
				IsFolder: false,
				Name:     entry.Name(),
			})
		}
	}
	return res, nil
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
