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
	err = vmstorage.WriteJson(fs.GetManifestPath(), manifest)
	return err
}

// CreateDisk returns the name of a temporary file where it is possible to store chunks
func (fs *FileSystemWrapper) CreateDisk(diskName string) (string, error) {
	var err error = vmstorage.CreateFolderRec(fs.basePath)
	if err != nil {
		return "", err
	}
	randomString, err := utils.RandomString(16)
	if err != nil {
		return "", err
	}
	var tempFileName string = fmt.Sprintf("%s_%s.tmp", randomString, diskName)
	err = vmstorage.CreateFile(fs.GetDiskPath(tempFileName))
	if err != nil {
		return "", err
	}
	return tempFileName, nil
}

func (fs *FileSystemWrapper) WriteChunk(tempDiskName string, byteIndex int64, chunk io.Reader) error {
	return vmstorage.WriteFileChunk(fs.GetDiskPath(tempDiskName), byteIndex, chunk)
}

func (fs *FileSystemWrapper) CommitDisk(tempDiskName string, diskName string) error {
	var randomString string = strings.Split(tempDiskName, "_")[0]
	if fmt.Sprintf("%s_%s.tmp", randomString, diskName) != tempDiskName {
		return errors.New("disk name and temporary disk name do not match")
	}
	var err error = vmstorage.RenameFile(fs.GetDiskPath(tempDiskName), fs.GetDiskPath(diskName))
	if err != nil {
		return err
	}
	err = vmstorage.DeleteFile(tempDiskName)
	return err
}
