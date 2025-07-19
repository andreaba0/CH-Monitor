package vmmanager

import vmstorage "vmm/storage"

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
}

func NewHypervisorMonitor(fs vmstorage.FileSystemStorage) *HypervisorMonitor {
	return &HypervisorMonitor{
		fs: fs,
	}
}

func (hm *HypervisorMonitor) LoadVirtualMachines(runningInstances []RunningCHInstance, manifestList []vmstorage.Manifest) error {
	var i int
	var err error
	for i = 0; i < len(manifestList); i++ {
		var manifest = manifestList[i]
		var vmId = manifest.GuestName
		hm.VirtualMachines[vmId] = VirtualMachine{
			PID:      nil,
			State:    Unknow,
			Manifest: &manifest,
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
			continue
		}
		if vmId != hm.VirtualMachines[vmId].Manifest.GuestName {
			continue
		}
		hm.VirtualMachines[vmId] = VirtualMachine{
			PID:      &instance.PID,
			Manifest: hm.VirtualMachines[vmId].Manifest,
			State:    hm.VirtualMachines[vmId].State,
		}
	}
	return nil
}

type HypervisorMonitorService interface {
	LoadVirtualMachines(runningInstances []RunningCHInstance, manifestList []vmstorage.Manifest) error
}
