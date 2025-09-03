package virtualmachine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"vmm/utils"

	"go.uber.org/zap"
)

type FileSystemWrapper struct {
	basePath string
	logger   *zap.Logger
}

func (fs *FileSystemWrapper) GetManifestPath() string {
	return filepath.Join(fs.basePath, "manifest.json")
}

func (fs *FileSystemWrapper) GetKernelPath(kernelName string) string {
	return filepath.Join(fs.basePath, "kernel", kernelName)
}

func (fs *FileSystemWrapper) GetKernelStoragePath() string {
	return filepath.Join(fs.basePath, "kernel")
}

func (fs *FileSystemWrapper) GetDiskPath(diskName string) string {
	return filepath.Join(fs.basePath, "disks", diskName)
}

func (fs *FileSystemWrapper) GetDiskStoragePath() string {
	return filepath.Join(fs.basePath, "disks")
}

func (fs *FileSystemWrapper) ReadManifest() (*Manifest, error) {
	fileBytes, err := os.ReadFile(fs.GetManifestPath())
	if err != nil {
		return nil, err
	}
	var manifest *Manifest
	err = json.Unmarshal(fileBytes, &manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (fs *FileSystemWrapper) StoreManifest(manifest *Manifest) error {
	var err error = fs.createFolderRecursively(fs.basePath)
	if err != nil {
		return err
	}
	err = fs.createFile(fs.GetManifestPath())
	if err != nil {
		return err
	}
	content, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	return os.WriteFile(fs.GetManifestPath(), content, os.ModePerm)
}

func (fs *FileSystemWrapper) createFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	return nil
}

func (fs *FileSystemWrapper) createTempFile(storagePath string, fileName string) (string, error) {
	var err error = fs.createFolderRecursively(storagePath)
	if err != nil {
		return "", err
	}

	randomString, err := utils.RandomString(16)
	if err != nil {
		return "", err
	}

	var tmpFileName string = fmt.Sprintf("%s_%s.tmp", randomString, fileName)
	err = fs.createFile(filepath.Join(storagePath, tmpFileName))
	if err != nil {
		return "", err
	}
	return tmpFileName, nil
}

func (fs *FileSystemWrapper) createFolderRecursively(folderPath string) error {
	return os.MkdirAll(folderPath, os.ModePerm)
}

func (fs *FileSystemWrapper) CreateDisk(diskName string) (string, error) {
	return fs.createTempFile(fs.GetDiskStoragePath(), diskName)
}

func (fs *FileSystemWrapper) CreateKernel(kernelName string) (string, error) {
	return fs.createTempFile(fs.GetKernelStoragePath(), kernelName)
}

func (fs *FileSystemWrapper) writeChunk(fileFullPath string, byteIndex int64, chunk io.Reader) error {
	fd, err := os.OpenFile(fileFullPath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Seek(byteIndex, 0)
	if err != nil {
		return err
	}
	_, err = io.Copy(fd, chunk)
	return err
}

func (fs *FileSystemWrapper) WriteDiskChunk(tmpDiskName string, byteIndex int64, chunk io.Reader) error {
	return fs.writeChunk(fs.GetDiskPath(tmpDiskName), byteIndex, chunk)
}

func (fs *FileSystemWrapper) WriteKernelChunk(tmpKernelName string, byteIndex int64, chunk io.Reader) error {
	return fs.writeChunk(fs.GetKernelPath(tmpKernelName), byteIndex, chunk)
}

func (fs *FileSystemWrapper) commitOperation(storagePath string, tmpFileName string, fileName string) error {
	var randomString string = strings.Split(tmpFileName, "_")[0]
	if fmt.Sprintf("%s_%s.tmp", randomString, fileName) != tmpFileName {
		return errors.New("disk name and temporary disk name do not match")
	}
	return os.Rename(filepath.Join(storagePath, tmpFileName), filepath.Join(storagePath, fileName))
}

func (fs *FileSystemWrapper) CommitDisk(tmpDiskName string, diskName string) error {
	return fs.commitOperation(fs.GetDiskStoragePath(), tmpDiskName, diskName)
}

func (fs *FileSystemWrapper) CommitKernel(tmpKernelName string, kernelName string) error {
	return fs.commitOperation(fs.GetKernelStoragePath(), tmpKernelName, kernelName)
}
