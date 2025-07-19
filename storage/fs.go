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

func (fs *FileSystemStorage) GetManifestPath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "manifest.json")
}

func (fs *FileSystemStorage) GetDiskPath(vmId string, diskName string) string {
	return filepath.Join(fs.basePath, vmId, "disks", diskName)
}

func (fs *FileSystemStorage) GetSocketPath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "cloud-hypervisor-vm.sock")
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
	var manifest *Manifest = &Manifest{}
	var content []byte
	content, err = os.ReadFile(fs.GetManifestPath(vmId))
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
	err = os.MkdirAll(fs.GetDiskStoragePath(vmId), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(fs.GetManifestPath(vmId), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileSystemStorage) GetVirtualMachineIdFromSocket(socketFilePath string) (string, error) {
	rel, err := filepath.Rel(fs.basePath, socketFilePath)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", errors.New("invalid socket path")
	}
	parts := strings.Split(rel, string(os.PathSeparator))
	if len(parts) < 2 {
		return "", errors.New("unexpected socket structure")
	}
	return parts[0], nil
}

func (fs *FileSystemStorage) GetFullVirtualMachineList() ([]Manifest, error) {
	var res []Manifest
	folders, err := os.ReadDir(fs.basePath)
	if err != nil {
		return []Manifest{}, errors.New("unable to read from folder")
	}
	for i := 0; i < len(folders); i++ {
		content, err := os.ReadFile(fs.GetManifestPath(folders[i].Name()))
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
