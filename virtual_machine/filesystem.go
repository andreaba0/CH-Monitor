package virtualmachine

import (
	"io"
	"path/filepath"
	vmstorage "vmm/storage"

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

func (fs *FileSystemWrapper) CreateDisk(diskName string) error {
	var err error = vmstorage.CreateFolderRec(fs.basePath)
	if err != nil {
		return err
	}
	return vmstorage.CreateFile(fs.GetDiskPath(diskName))
}

func (fs *FileSystemWrapper) WriteChunk(diskName string, byteIndex int64, chunk io.Reader) error {
	return vmstorage.WriteFileChunk(fs.GetDiskPath(diskName), byteIndex, chunk)
}
