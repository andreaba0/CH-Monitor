package virtualmachine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"
	"vmm/metadata"

	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

type VirtualMachine struct {
	manifest        *Manifest
	hypervisor      *cloudhypervisor.CloudHypervisor
	storage         *FileSystemWrapper
	logger          *zap.Logger
	mu              sync.Mutex
	metadataManager *metadata.MetadataManager
	defaultBridge   netlink.Link
}

func NewVirtualMachine(manifest *Manifest, logger *zap.Logger, storagePath string, defaultBridge string, metadataService *metadata.MetadataManager) (*VirtualMachine, error) {
	bridgeLink, err := netlink.LinkByName(defaultBridge)
	if err != nil {
		return nil, err
	}
	return &VirtualMachine{
		manifest:   manifest,
		hypervisor: nil,
		storage: &FileSystemWrapper{
			basePath: filepath.Join(storagePath, manifest.GuestIdentifier.String()),
			logger:   logger,
		},
		logger:          logger,
		metadataManager: metadataService,
		defaultBridge:   bridgeLink,
	}, nil

}

func LoadVirtualMachine(vmFolder string, logger *zap.Logger, defaultBridge string, metadataService *metadata.MetadataManager) (*VirtualMachine, error) {
	bridgeLink, err := netlink.LinkByName(defaultBridge)
	if err != nil {
		return nil, err
	}
	var storage *FileSystemWrapper = &FileSystemWrapper{
		basePath: vmFolder,
		logger:   logger,
	}
	manifest, err := storage.ReadManifest()
	if err != nil {
		return nil, errors.New("unable to read manifest")
	}
	return &VirtualMachine{
		manifest:        manifest,
		hypervisor:      nil,
		storage:         storage,
		logger:          logger,
		metadataManager: metadataService,
		defaultBridge:   bridgeLink,
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

func (vm *VirtualMachine) RunInstance(binaryPath string, remoteUri string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.hypervisor != nil {
		return errors.New("an instance is already running")
	}
	cloudhypervisor, err := cloudhypervisor.NewCloudHypervisor(binaryPath, remoteUri)
	if err != nil {
		return err
	}
	vm.AttachInstance(cloudhypervisor)
	return nil
}

func (vm *VirtualMachine) RequestBoot(binaryPath string, remoteUri string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.hypervisor != nil {
		return errors.New("virtual machine is already running")
	}
	hypervisor, err := cloudhypervisor.NewCloudHypervisor(binaryPath, remoteUri)
	if err != nil {
		return err
	}
	vm.hypervisor = hypervisor
	err = vm.createVirtualMachine()
	if err != nil {
		return err
	}
	err = vm.bootVirtualMachine()
	if err != nil {
		return err
	}
	err = vm.connectNetworking()
	if err != nil {
		return errors.New("there was an error connecting vm to network interfaces")
	}
	return nil
}

func (vm *VirtualMachine) RequestShutdown() error {
	err := vm.shutdownVirtualMachine()
	if err != nil {
		return err
	}
	return nil
}

func (vm *VirtualMachine) createVirtualMachine() error {
	vmManifest, err := vm.parseManifestToCloudHypervisor()
	if err != nil {
		return err
	}
	arr, err := json.Marshal(vmManifest)
	if err != nil {
		return err
	} else {
		fmt.Print(string(arr))
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(vmManifest)
	if err != nil {
		return err
	}
	uri, _ := vm.hypervisor.RestServer.GetUri(cloudhypervisor.CREATE)
	req, err := http.NewRequest(http.MethodPut, uri, &buf)
	if err != nil {
		return err
	}
	res, err := vm.hypervisor.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}
	return nil
}

func (vm *VirtualMachine) bootVirtualMachine() error {
	uri, _ := vm.hypervisor.RestServer.GetUri(cloudhypervisor.BOOT)
	req, err := http.NewRequest(http.MethodPut, uri, nil)
	if err != nil {
		return err
	}
	res, err := vm.hypervisor.HttpClient.Do(req)
	if err != nil || (res.StatusCode < 200 || res.StatusCode > 299) {
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	defer res.Body.Close()
	return nil
}

func (vm *VirtualMachine) shutdownVirtualMachine() error {
	uri, _ := vm.hypervisor.RestServer.GetUri(cloudhypervisor.SHUTDOWN)
	req, err := http.NewRequest(http.MethodPut, uri, nil)
	if err != nil {
		return err
	}
	res, err := vm.hypervisor.HttpClient.Do(req)
	if err != nil || (res.StatusCode < 200 || res.StatusCode > 299) {
		return errors.New("there was an error performing http request")
	}
	return nil
}

func (vm *VirtualMachine) connectNetworking() error {
	/*for i := 0; i < len(vm.manifest.Config.Networks); i++ {
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
	}*/
	return nil
}

func (vm *VirtualMachine) parseManifestToCloudHypervisor() (*cloudhypervisor.Manifest, error) {
	chManifest := &cloudhypervisor.Manifest{
		Cpus: cloudhypervisor.VmCpus{
			Boot_vcpus: vm.manifest.Config.Cpus,
			Max_vcpus:  vm.manifest.Config.Cpus,
		},
		Platform: cloudhypervisor.Platform{
			Uuid: vm.manifest.GuestIdentifier.String(),
		},
		Rng: cloudhypervisor.Rng{
			Src: vm.manifest.Config.Rng.Src,
		},
		Serial: cloudhypervisor.Serial{
			Mode: "File",
			File: fmt.Sprintf("/tmp/%s.log", vm.manifest.GuestIdentifier.String()),
		},
		Console: cloudhypervisor.Console{
			Mode: "Off",
		},
	}
	disks := []cloudhypervisor.Disk{}
	for i := 0; i < len(vm.manifest.Config.Disks); i++ {
		disks = append(disks, cloudhypervisor.Disk{
			Path: vm.storage.GetDiskPath(vm.manifest.Config.Disks[i].Name),
		})
	}
	chManifest.Disks = disks
	if vm.manifest.Config.Kernel != "" && vm.manifest.Config.Init != "" {
		chManifest.Payload = cloudhypervisor.Payload{
			Kernel:  vm.storage.GetKernelPath(vm.manifest.Config.Kernel),
			Cmdline: fmt.Sprintf("console=ttyS0 root=/dev/vda rw init=%s", vm.manifest.Config.Init),
		}
	}
	return chManifest, nil
}
