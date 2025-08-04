package virtualmachine

import (
	"errors"
	"io"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"

	"go.uber.org/zap"
)

type VirtualMachine struct {
	manifest   *Manifest
	hypervisor *cloudhypervisor.CloudHypervisor
	storage    *FileSystemWrapper
	logger     *zap.Logger
	mu         sync.Mutex
}

func NewVirtualMachine(manifest *Manifest, logger *zap.Logger, vmPath string) *VirtualMachine {
	return &VirtualMachine{
		manifest:   manifest,
		hypervisor: nil,
		storage: &FileSystemWrapper{
			basePath: vmPath,
			logger:   logger,
		},
		logger: logger,
	}

}

func LoadVirtualMachine(vmFolder string, logger *zap.Logger) (*VirtualMachine, error) {
	var storage *FileSystemWrapper = &FileSystemWrapper{
		basePath: vmFolder,
		logger:   logger,
	}
	var err error
	var manifest *Manifest
	manifest, err = storage.ReadManifest()
	if err != nil {
		return nil, errors.New("unable to read manifest")
	}
	return &VirtualMachine{
		manifest:   manifest,
		hypervisor: nil,
		storage:    storage,
		logger:     logger,
	}, nil
}

func (vm *VirtualMachine) StoreManifest() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.StoreManifest(vm.manifest)
}

func (vm *VirtualMachine) GetManifest() *Manifest {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.manifest
}

func (vm *VirtualMachine) CreateDisk(diskName string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.CreateDisk(diskName)
}

func (vm *VirtualMachine) WriteChunkToDisk(diskName string, byteIndex int64, chunk io.Reader) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.WriteChunk(diskName, byteIndex, chunk)
}

func (vm *VirtualMachine) AttachInstance(hypervisor *cloudhypervisor.CloudHypervisor) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.hypervisor = hypervisor
}
