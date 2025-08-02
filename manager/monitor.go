package vmmanager

import (
	"sync"
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
	VirtualMachines map[string]virtualmachine.VirtualMachine
	vmsMu           sync.Mutex
	fs              *virtualmachine.FileSystemWrapper
	logger          *zap.Logger
}

func NewHypervisorMonitor(fs *virtualmachine.FileSystemWrapper, logger *zap.Logger) *HypervisorMonitor {
	return &HypervisorMonitor{
		fs:     fs,
		logger: logger,
	}
}

func (hm *HypervisorMonitor) LoadVirtualMachines(runningInstances []RunningCHInstance, manifestList []*virtualmachine.Manifest) error {
	var i int
	var err error
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
	}
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
