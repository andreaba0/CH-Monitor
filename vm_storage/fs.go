package vmstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemStorage struct {
	basePath string
}

func NewFileSystemStorage(path string) (*FileSystemStorage, error) {
	if path == "" {
		return nil, errors.New("path to vm storage required")
	}

	return &FileSystemStorage{
		basePath: path,
	}, nil
}

func (fs *FileSystemStorage) CreateDisk(vmId string, fileName string, diskSize int64) error {
	var err error
	var fd *os.File
	var tempFileName string
	var location string = fs.GetDiskStoragePath(vmId)

	tempFileName = fmt.Sprintf("%s.tmp", fileName)

	fd, err = os.OpenFile(filepath.Join(location, tempFileName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = fd.Truncate(diskSize)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileSystemStorage) WriteDiskChunk(vmId string, fileName string, byteIndex int64, chunk io.Reader) error {
	var err error
	var fd *os.File

	var fullFilePath = filepath.Join(fs.GetDiskStoragePath(vmId), fmt.Sprintf("%s.tmp", fileName))

	fd, err = os.OpenFile(fullFilePath, os.O_WRONLY, os.ModePerm)
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

func (fs *FileSystemStorage) CompleteDiskWrite(vmId string, fileName string) error {
	var err error
	var location string = fs.GetDiskStoragePath(vmId)
	var oldFilename string = filepath.Join(location, fmt.Sprintf("%s.tmp", fileName))
	var newFilename string = filepath.Join(location, fileName)
	err = os.Rename(oldFilename, newFilename)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystemStorage) GetDiskStoragePath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "disks")
}

func (fs *FileSystemStorage) ReadManifest(vmId string) (*Manifest, error) {
	var err error
	var manifest *Manifest
	var content []byte
	content, err = os.ReadFile(filepath.Join(fs.basePath, vmId, "manifest.json"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (fs *FileSystemStorage) CreateVirtualMachine(vmId string, manifest *Manifest) error {
	var err error
	var content []byte

	content, err = json.Marshal(manifest)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(fs.basePath, vmId), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(fs.basePath, vmId, "manifest.json"), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileSystemStorage) GetIdFromSocket(socketPath string) (string, error) {
	var parts = strings.Split(socketPath, "/")
	var basePathParts = strings.Split(fs.basePath, "/")
	var i int
	if len(parts) != len(basePathParts) || len(parts) == len(basePathParts) || len(parts) != len(basePathParts)+2 {
		return "", errors.New("socket is expected to be in another folder")
	}
	for i = 0; i < len(basePathParts); i++ {
		if parts[i] != basePathParts[i] {
			return "", errors.New("socket is expected to be in another folder")
		}
	}
	var vmId = parts[len(parts)-2]
	return vmId, nil
}

func (fs *FileSystemStorage) GetFullVirtualMachineList() ([]Manifest, error) {
	var res []Manifest
	folders, err := os.ReadDir(fs.basePath)
	if err != nil {
		return []Manifest{}, errors.New("unable to read from folder")
	}
	for i := 0; i < len(folders); i++ {
		content, err := os.ReadFile(filepath.Join(fs.basePath, folders[i].Name(), "manifest.json"))
		if err != nil {
			continue
		}
		var manifest Manifest = Manifest{}
		err = json.Unmarshal(content, &manifest)
		if err != nil {
			continue
		}
		res = append(res, manifest)
	}
	return res, nil
}

func (fs *FileSystemStorage) CreateManifest(vmId string, manifest Manifest) error {
	return nil
}

type FileSystemStorageService interface {
	CreateDisk(vmId string, fileName string, diskSize int64) error
	WriteDiskChunk(vmId string, fileName string, byteIndex int64, chunk io.Reader) error
	CompleteDiskWrite(vmId string, fileName string) error
	CreateManifest(vmId string, manifest Manifest) error
	ReadManifest(vmId string) (*Manifest, error)
	//Delete(vmId string) error
	GetDiskStoragePath(vmId string) string
	GetFullVirtualMachineList() ([]Manifest, error)
	CreateVirtualMachine(vmId string, manifest *Manifest) error
	GetIdFromSocket(socketPath string) (string, error)
}
