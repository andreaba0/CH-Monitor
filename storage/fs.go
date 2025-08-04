package vmstorage

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func CreateFile(path string) error {
	var err error
	var fd *os.File

	fd, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func CreateFileIfNotExists(filePath string) error {
	var err error
	var fd *os.File

	fd, err = os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func WriteFileChunk(path string, byteIndex int64, chunk io.Reader) error {
	var err error
	var fd *os.File
	fd, err = os.OpenFile(path, os.O_WRONLY, os.ModePerm)
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

func RenameFile(oldPath string, newPath string) error {
	var err error = os.Rename(oldPath, newPath)
	return err
}

func ReadYaml[T any](path string) (T, error) {
	var res T
	var err error
	var dataByte []byte = []byte{}
	dataByte, err = os.ReadFile(path)
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
	var err error
	var dataByte []byte = []byte{}
	dataByte, err = os.ReadFile(path)
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
	var err error
	var dataByte []byte
	var fd *os.File
	dataByte, err = json.Marshal(content)
	if err != nil {
		return err
	}
	fd, err = os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(dataByte)
	return err
}

func CreateFolderRec(path string) error {
	var err error = os.MkdirAll(path, os.ModePerm)
	return err
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

/*
func GetFullVirtualMachineList[T any]() ([]T, error) {
	var res []T
	folders, err := os.ReadDir(fs.basePath)
	if err != nil {
		return []FileContent{}, errors.New("unable to read from folder")
	}
	for i := 0; i < len(folders); i++ {
		content, err := os.ReadFile(fs.GetManifestPath(folders[i].Name()))
		if err != nil {
			continue
		}
		res = append(res, FileContent{
			Path:    fs.GetManifestPath(folders[i].Name()),
			Content: content,
		})
	}
	return res, nil
}
*/
