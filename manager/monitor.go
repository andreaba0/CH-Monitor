package vmmanager

import (
	"encoding/json"
	vmstorage "vmm/storage"

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
	VirtualMachines map[string]VirtualMachine
	fs              vmstorage.FileSystemStorage
	logger          *zap.Logger
}

func NewHypervisorMonitor(fs vmstorage.FileSystemStorage, logger *zap.Logger) *HypervisorMonitor {
	return &HypervisorMonitor{
		fs:     fs,
		logger: logger,
	}
}

func (hm *HypervisorMonitor) LoadVirtualMachines(runningInstances []RunningCHInstance, manifestContentList []vmstorage.FileContent) error {
	var i int
	var err error
	for i = 0; i < len(manifestContentList); i++ {
		var manifest *Manifest = &Manifest{}
		err = json.Unmarshal(manifestContentList[i].Content, manifest)
		if err != nil {
			hm.logger.Error("Unable to parse file content. Maybe a corrupted file?", zap.String("path", manifestContentList[i].Path))
			return err
		}
		var vmId = manifest.GuestName
		hm.VirtualMachines[vmId] = VirtualMachine{
			PID:      nil,
			Manifest: manifest,
		}
	}
	for i = 0; i < len(runningInstances); i++ {
		var instance RunningCHInstance = runningInstances[i]
		var vmId string
		vmId, err = hm.fs.GetVirtualMachineIdFromSocket(instance.UnixSocketPath)
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
