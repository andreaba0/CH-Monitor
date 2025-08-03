package vmm

import "fmt"

type VirtualMachineAction int

const (
	CREATE VirtualMachineAction = iota
	DELETE
	BOOT
	SHUTDOWN
)

type HypervisorBinary struct {
	BinaryPath string
	RemoteUri  string
}

func (hb *HypervisorBinary) GetLaunchCommandString(socket string) string {
	return fmt.Sprintf("%s --api-socket path=%s", hb.BinaryPath, socket)
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
