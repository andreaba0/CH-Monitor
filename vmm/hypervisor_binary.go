package vmm

import "fmt"

type VirtualMachineAction int

const (
	CREATE VirtualMachineAction = iota
	DELETE
	BOOT
	SHUTDOWN
	INFO
)

type HypervisorBinary struct {
	remoteUri string
}

func NewHypervisorBinary(remoteUri string) *HypervisorBinary {
	return &HypervisorBinary{
		remoteUri: remoteUri,
	}
}

func (hb *HypervisorBinary) GetUri(action VirtualMachineAction) *string {
	switch action {
	case CREATE:
		str := fmt.Sprintf("%s/vm.create", hb.remoteUri)
		return &str
	case BOOT:
		str := fmt.Sprintf("%s/vm.boot", hb.remoteUri)
		return &str
	case DELETE:
		str := fmt.Sprintf("%s/vm.delete", hb.remoteUri)
		return &str
	case SHUTDOWN:
		str := fmt.Sprintf("%s/vm.shutdown", hb.remoteUri)
		return &str
	case INFO:
		str := fmt.Sprintf("%s/vm.info", hb.remoteUri)
		return &str
	default:
		return nil
	}
}
