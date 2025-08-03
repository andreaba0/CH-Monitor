package vmm

import (
	"path/filepath"
	"sync"
	vmstorage "vmm/storage"
	virtualmachine "vmm/virtual_machine"

	"go.uber.org/zap"
)

type VMConfig struct {
	RootfsName    string
	AttachedDisks []string
	Kernel        string
	UnixSocket    string
	TapDevice     string
	PID           int
}

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

	/*var i int
	for i = 0; i < len(manifestList); i++ {
		var vmId = manifestList[i].GuestIdentifier
		hm.VirtualMachines[vmId] = VirtualMachine{
			PID:      nil,
			Manifest: manifestList[i],
		}
	}
	for i = 0; i < len(runningInstances); i++ {
		var instance RunningCHInstance = runningInstances[i]
		var vmId string
		vmId, err = instance.GetVirtualMachineIdFromSocket(hm.fs.basePath)
		if err != nil {
			continue
		}
		var ok bool
		_, ok = hm.VirtualMachines[vmId]
		if !ok {
			hm.logger.Error("There is an instance running without a manifest", zap.String("path", hm.fs.GetManifestPath(instance.UnixSocketPath)))
			continue
		}
		if vmId != hm.VirtualMachines[vmId].Manifest.GuestName {
			continue
		}
		hm.VirtualMachines[vmId] = VirtualMachine{
			PID:      &instance.PID,
			Manifest: hm.VirtualMachines[vmId].Manifest,
		}
	}*/
	return nil
}

func (hm *HypervisorMonitor) MergeRunningInstances() error {
	return nil
}

/*func GetVirtualMachineList() ([]*Manifest, error) {
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
}*/
