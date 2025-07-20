package vmmanager

import vmstorage "vmm/storage"

type VirtualMachineState int

const (
	Created VirtualMachineState = iota
	Booted
	Running
	Paused
	Unknow
)

type VirtualMachine struct {
	Manifest          *Manifest
	PID               *int
	fileSystemStorage *vmstorage.FileSystemStorage
}

func (vm *VirtualMachine) Create() {

}

func (vm *VirtualMachine) Boot() {
	var _ string = vm.fileSystemStorage.GetSocketPath(vm.Manifest.GuestName)
}

func (vm *VirtualMachine) Shutdown() {

}

func (vm *VirtualMachine) GetState() VirtualMachineState {
	return VirtualMachineState(Unknow)
}
