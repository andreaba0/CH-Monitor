package virtualmachine

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	vmstorage "vmm/storage"
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
	manifest, err := vmstorage.ReadJson[*Manifest](fs.GetManifestPath())
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (fs *FileSystemWrapper) StoreManifest(manifest *Manifest) error {
	var err error = vmstorage.CreateFolderRec(fs.basePath)
	if err != nil {
		return err
	}
	err = vmstorage.CreateFileIfNotExists(fs.GetManifestPath())
	if err != nil {
		return err
	}
	err = vmstorage.WriteJson(fs.GetManifestPath(), manifest)
	return err
}

func (fs *FileSystemWrapper) createFile(storagePath string, fileName string) (string, error) {
	var err error = vmstorage.CreateFolderRec(storagePath)
	if err != nil {
		return "", err
	}

	randomString, err := utils.RandomString(16)
	if err != nil {
		return "", err
	}

	var tmpFileName string = fmt.Sprintf("%s_%s.tmp", randomString, fileName)
	err = vmstorage.CreateFile(filepath.Join(storagePath, tmpFileName))
	if err != nil {
		return "", err
	}
	return tmpFileName, nil
}

func (fs *FileSystemWrapper) CreateDisk(diskName string) (string, error) {
	return fs.createFile(fs.GetDiskStoragePath(), diskName)
}

func (fs *FileSystemWrapper) CreateKernel(kernelName string) (string, error) {
	return fs.createFile(fs.GetKernelStoragePath(), kernelName)
}

func (fs *FileSystemWrapper) writeChunk(fileFullPath string, byteIndex int64, chunk io.Reader) error {
	return vmstorage.WriteFileChunk(fileFullPath, byteIndex, chunk)
}

func (fs *FileSystemWrapper) WriteDiskChunk(tmpDiskName string, byteIndex int64, chunk io.Reader) error {
	return fs.writeChunk(fs.GetDiskPath(tmpDiskName), byteIndex, chunk)
}

func (fs *FileSystemWrapper) WriteKernelChunk(tmpKernelName string, byteIndex int64, chunk io.Reader) error {
	return fs.writeChunk(fs.GetKernelPath(tmpKernelName), byteIndex, chunk)
}

func (fs *FileSystemWrapper) commitOperation(storagePath string, tmpFileName string, fileName string) error {
	var err error = vmstorage.CreateFolderRec(storagePath)
	if err != nil {
		return err
	}
	var randomString string = strings.Split(tmpFileName, "_")[0]
	if fmt.Sprintf("%s_%s.tmp", randomString, fileName) != tmpFileName {
		return errors.New("disk name and temporary disk name do not match")
	}
	err = vmstorage.RenameFile(filepath.Join(storagePath, tmpFileName), filepath.Join(storagePath, fileName))
	return err
}

func (fs *FileSystemWrapper) CommitDisk(tmpDiskName string, diskName string) error {
	return fs.commitOperation(fs.GetDiskStoragePath(), tmpDiskName, diskName)
}

func (fs *FileSystemWrapper) CommitKernel(tmpKernelName string, kernelName string) error {
	return fs.commitOperation(fs.GetKernelStoragePath(), tmpKernelName, kernelName)
}
