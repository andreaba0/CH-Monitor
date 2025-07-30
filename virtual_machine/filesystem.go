package virtualmachine

import (
	"errors"
	"fmt"
	"path/filepath"
	vmstorage "vmm/storage"

	"go.uber.org/zap"
)

type FileSystemWrapper struct {
	basePath string
	logger   *zap.Logger
}

func NewFileSystemWrapper(path string, logger *zap.Logger) (*FileSystemWrapper, error) {
	if path == "" {
		return nil, errors.New("path to vm storage required")
	}

	return &FileSystemWrapper{
		basePath: path,
		logger:   logger,
	}, nil
}

func (fs *FileSystemWrapper) GetManifestPath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "manifest.json")
}

func (fs *FileSystemWrapper) GetDiskPath(vmId string, diskName string) string {
	return filepath.Join(fs.basePath, vmId, "disks", diskName)
}

func (fs *FileSystemWrapper) GetSocketPath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "cloud-hypervisor-vm.sock")
}

func (fs *FileSystemWrapper) GetDiskStoragePath(vmId string) string {
	return filepath.Join(fs.basePath, vmId, "disks")
}

func (fs *FileSystemWrapper) CreateVirtualMachine(vmId string, manifest *Manifest) error {
	var err error = vmstorage.CreateFolderRec(fs.GetDiskStoragePath(vmId))
	if err != nil {
		return err
	}
	err = vmstorage.WriteJson(fs.GetManifestPath(vmId), manifest)
	for i := 0; i < len(manifest.Config.Disks); i++ {
		var tmpFileName string = fmt.Sprintf("%s.tmp", manifest.Config.Disks[i].Name)
		err = vmstorage.CreateFile(fs.GetDiskPath(vmId, tmpFileName))
		if err != nil {
			fs.logger.Error("Unable to create file", zap.String("path", fs.GetDiskPath(vmId, tmpFileName)))
		}
	}
	return err
}

func (fs *FileSystemWrapper) ReadManifest(vmId string) (*Manifest, error) {
	manifest, err := vmstorage.ReadJson[*Manifest](fs.GetManifestPath(vmId))
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func (fs *FileSystemWrapper) GetVirtualMachineList() ([]*Manifest, error) {
	var res []*Manifest = []*Manifest{}
	entries, err := vmstorage.ListFolder(fs.basePath)
	if err != nil {
		return []*Manifest{}, err
	}
	for _, entry := range entries {
		if !entry.IsFolder {
			continue
		}
		manifest, err := vmstorage.ReadJson[*Manifest](fs.GetManifestPath(entry.Name))
		if err != nil {
			fs.logger.Error("Unable to read manifest from file", zap.String("path", fs.GetManifestPath(entry.Name)))
		}
		res = append(res, manifest)
	}
	return res, nil
}
