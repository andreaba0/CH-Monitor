package vmmanager

import vmstorage "vmm/vm_storage"

type VirtualMachineState int

const (
	Created VirtualMachineState = iota
	Booted
	Running
	Paused
	Unknow
)

type VirtualMachine struct {
	State    VirtualMachineState
	Manifest *vmstorage.Manifest
	PID      *int
}
