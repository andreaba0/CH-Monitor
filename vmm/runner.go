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
	VirtualMachines map[string]*virtualmachine.VirtualMachine
	vmsMu           sync.Mutex
	logger          *zap.Logger
}

func NewHypervisorMonitor(logger *zap.Logger) *HypervisorMonitor {
	return &HypervisorMonitor{
		VirtualMachines: make(map[string]*virtualmachine.VirtualMachine),
		logger:          logger,
	}
}

func MonitorSetup(manifestPath string, vmm *HypervisorMonitor) error {
	var manifest Manifest
	var err error
	manifest, err = vmstorage.ReadYaml[Manifest](manifestPath)
	if err != nil {
		return err
	}
	var binaryPath string = manifest.HypervisorPath
	var remoteUri string = manifest.HypervisorSocketUri
	var hypervisorBinary HypervisorBinary = *NewHypervisorBinary(remoteUri)
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
		vm, err := virtualmachine.LoadVirtualMachine(filepath.Join(basePath, entry.Name), hm.logger)
		if err != nil {
			hm.logger.Error("Unable to read manifest from file", zap.String("base_path", basePath), zap.String("vm_id", entry.Name))
		}
		guestName := vm.GetManifest().GuestIdentifier
		hm.VirtualMachines[guestName.String()] = vm
	}
	return nil
}

func (hm *HypervisorMonitor) MergeRunningInstances(hypervisorBinaryPath string, hypervisorBinary *HypervisorBinary) error {
	var err error
	var instances []*cloudhypervisor.CloudHypervisor
	instances, err = LoadProcessData(hypervisorBinaryPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(instances); i++ {
		var manifest cloudhypervisor.Manifest
		var err error
		res, err := instances[i].HttpClient.Get(*hypervisorBinary.GetUri(VirtualMachineAction(INFO)))
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
		vm = hm.VirtualMachines[manifest.Platform.Uuid.String()]
		hm.vmsMu.Unlock()
		vm.AttachInstance(instances[i])
	}
	return nil
}
