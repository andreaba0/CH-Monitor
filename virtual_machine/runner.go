package virtualmachine

import (
	"errors"
	"io"
	"net"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"
	vmnetworking "vmm/vm_networking"

	"go.uber.org/zap"
)

type VirtualMachine struct {
	manifest       *Manifest
	hypervisor     *cloudhypervisor.CloudHypervisor
	storage        *FileSystemWrapper
	logger         *zap.Logger
	mu             sync.Mutex
	networkManager *vmnetworking.NetworkManager
}

func NewVirtualMachine(manifest *Manifest, logger *zap.Logger, vmPath string, defaultBridge string) (*VirtualMachine, error) {
	nm, err := vmnetworking.NewNetworkManager(defaultBridge)
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{
		manifest:   manifest,
		hypervisor: nil,
		storage: &FileSystemWrapper{
			basePath: vmPath,
			logger:   logger,
		},
		logger:         logger,
		networkManager: nm,
	}, nil

}

func LoadVirtualMachine(vmFolder string, logger *zap.Logger, defaultBridge string) (*VirtualMachine, error) {
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
	nm, err := vmnetworking.NewNetworkManager(defaultBridge)
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{
		manifest:       manifest,
		hypervisor:     nil,
		storage:        storage,
		logger:         logger,
		networkManager: nm,
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

func (vm *VirtualMachine) CreateDisk(diskName string) (string, error) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.CreateDisk(diskName)
}

func (vm *VirtualMachine) CreateKernel(kernelName string) (string, error) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.CreateKernel(kernelName)
}

func (vm *VirtualMachine) WriteChunkToDisk(diskName string, byteIndex int64, chunk io.Reader) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.WriteDiskChunk(diskName, byteIndex, chunk)
}

func (vm *VirtualMachine) WriteChunkToKernel(kernelName string, byteIndex int64, chunk io.Reader) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.WriteKernelChunk(kernelName, byteIndex, chunk)
}

func (vm *VirtualMachine) CommitDisk(tempDiskName string, diskName string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.CommitDisk(tempDiskName, diskName)
}

func (vm *VirtualMachine) CommitKernel(tempKernelName string, kernelName string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	return vm.storage.CommitKernel(tempKernelName, kernelName)
}

func (vm *VirtualMachine) AttachInstance(hypervisor *cloudhypervisor.CloudHypervisor) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.hypervisor = hypervisor
}

func (vm *VirtualMachine) RunInstance() error {

	// Steps:
	// 1. Create cloud-hypervisor instance
	// 2. Connect tap interfaces to bridges
	// 3. Boot vm

	for i := 0; i < len(vm.manifest.Config.Networks); i++ {
		address := vm.manifest.Config.Networks[i]
		ip, ipNet, err := net.ParseCIDR(address.Address)
		if err != nil {
			return err
		}
		tap, err := vm.networkManager.GetTapInterface(ip, ipNet.Mask, vm.manifest.Tenant.String())
		if err != nil {
			return err
		}
		err = vm.networkManager.ConnectTapToDefault(tap)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(vm.manifest.Config.Vpc); i++ {
		address := vm.manifest.Config.Vpc[i]
		ip, ipNet, err := net.ParseCIDR(address.Address)
		if err != nil {
			return err
		}
		tap, err := vm.networkManager.GetTapInterface(ip, ipNet.Mask, vm.manifest.Tenant.String())
		if err != nil {
			return err
		}
		bridge, err := vm.networkManager.GetAndCreateIfNotExistsVpc(*ipNet, vm.manifest.Tenant.String())
		if err != nil {
			return err
		}
		err = vm.networkManager.ConnectTapToVpc(tap, bridge)
		if err != nil {
			return err
		}
	}
	return nil
}
