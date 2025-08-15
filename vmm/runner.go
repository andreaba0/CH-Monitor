package vmm

import (
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"
	vmstorage "vmm/storage"
	virtualmachine "vmm/virtual_machine"

	"go.uber.org/zap"
)

type HypervisorMonitor struct {
	virtualMachines map[string]*virtualmachine.VirtualMachine
	vmsMu           sync.Mutex
	logger          *zap.Logger
	manifest        *Manifest
}

func NewHypervisorMonitor(logger *zap.Logger) *HypervisorMonitor {
	return &HypervisorMonitor{
		virtualMachines: make(map[string]*virtualmachine.VirtualMachine),
		logger:          logger,
		manifest:        nil,
	}
}

func MonitorSetup(manifestPath string, vmm *HypervisorMonitor) error {
	var manifest Manifest
	var err error
	manifest, err = vmstorage.ReadYaml[Manifest](manifestPath)
	if err != nil {
		return err
	}
	vmm.SetManifest(&manifest)
	var binaryPath string = manifest.HypervisorPath
	var remoteUri string = manifest.HypervisorSocketUri
	var hypervisorBinary cloudhypervisor.HypervisorRestServer = *cloudhypervisor.NewHypervisorRestServer(remoteUri)
	err = vmm.LoadVirtualMachines(manifest.Server.StoragePath)
	if err != nil {
		return err
	}
	err = vmm.MergeRunningInstances(binaryPath, &hypervisorBinary)
	if err != nil {
		return err
	}
	return nil
}

func (hm *HypervisorMonitor) SetManifest(manifest *Manifest) {
	hm.manifest = manifest
}

func (hm *HypervisorMonitor) LoadVirtualMachines(basePath string) error {

	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	var err error
	entries, err := vmstorage.ListFolder(basePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsFolder {
			continue
		}
		vm, err := virtualmachine.LoadVirtualMachine(filepath.Join(basePath, entry.Name), hm.logger, hm.manifest.Bridge)
		if err != nil {
			hm.logger.Error("Unable to read manifest from file", zap.String("base_path", basePath), zap.String("vm_id", entry.Name))
		}
		guestName := vm.GetManifest().GuestIdentifier
		hm.virtualMachines[guestName.String()] = vm
	}
	return nil
}

func (hm *HypervisorMonitor) MergeRunningInstances(hypervisorBinaryPath string, hypervisorBinary *cloudhypervisor.HypervisorRestServer) error {
	var err error
	var instances []*cloudhypervisor.CloudHypervisor
	instances, err = LoadProcessData(hypervisorBinaryPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(instances); i++ {
		var manifest cloudhypervisor.Manifest
		var err error
		uri, err := hypervisorBinary.GetUri(cloudhypervisor.VirtualMachineAction(cloudhypervisor.INFO))
		if err != nil {
			return err
		}
		res, err := instances[i].HttpClient.Get(uri)
		if err != nil {
			return err
		}
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return errors.New("error while retrieving vm info")
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(resBody, &manifest)
		if err != nil {
			return err
		}
		var vm *virtualmachine.VirtualMachine
		hm.vmsMu.Lock()
		vm = hm.virtualMachines[manifest.Platform.Uuid]
		hm.vmsMu.Unlock()
		vm.AttachInstance(instances[i])
	}
	return nil
}

func (hm *HypervisorMonitor) CreateVirtualMachine(manifest *virtualmachine.Manifest) error {
	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	var err error
	vm, err := virtualmachine.NewVirtualMachine(manifest, hm.logger, filepath.Join(hm.manifest.Server.StoragePath, manifest.GuestIdentifier.String()), hm.manifest.Bridge)
	if err != nil {
		return err
	}
	err = vm.StoreManifest()
	if err != nil {
		return err
	}
	hm.virtualMachines[manifest.GuestIdentifier.String()] = vm
	return nil
}

func (hm *HypervisorMonitor) GetVirtualMachine(id string) *virtualmachine.VirtualMachine {
	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	return hm.virtualMachines[id]
}

func (hm *HypervisorMonitor) GetBinaryPath() string {
	return hm.manifest.HypervisorPath
}

func (hm *HypervisorMonitor) GetRestServerUri() string {
	return hm.manifest.HypervisorSocketUri
}
