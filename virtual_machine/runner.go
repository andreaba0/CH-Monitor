package virtualmachine

import (
	"fmt"
	cloudhypervisor "vmm/cloud_hypervisor"
)

type VirtualMachineAction int

const (
	CREATE VirtualMachineAction = iota
	DELETE
	BOOT
	SHUTDOWN
)

type HypervisorBinary struct {
	RemoteUri string
}

func (hb *HypervisorBinary) GetUri(action VirtualMachineAction) *string {
	switch action {
	case CREATE:
		str := fmt.Sprintf("%s/vm.create", hb.RemoteUri)
		return &str
	case BOOT:
		str := fmt.Sprintf("%s/vm.boot", hb.RemoteUri)
		return &str
	case DELETE:
		str := fmt.Sprintf("%s/vm.delete", hb.RemoteUri)
		return &str
	case SHUTDOWN:
		str := fmt.Sprintf("%s/vm.shutdown", hb.RemoteUri)
		return &str
	default:
		return nil
	}
}

type VirtualMachine struct {
	manifest   *Manifest
	hypervisor *cloudhypervisor.CloudHypervisor
	baseFolder string
	fs         *FileSystemWrapper
}

func NewVirtualMachine(manifest *Manifest, fs *FileSystemWrapper) (*VirtualMachine, error) {

}

func LoadStored(fs *FileSystemWrapper) []VirtualMachine {

}
