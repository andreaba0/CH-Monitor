package cloudhypervisor

import (
	"errors"
	"fmt"
)

type VirtualMachineAction int

const (
	CREATE VirtualMachineAction = iota
	DELETE
	BOOT
	SHUTDOWN
	INFO
)

type HypervisorRestServer struct {
	remoteUri string
}

func NewHypervisorRestServer(remoteUri string) *HypervisorRestServer {
	return &HypervisorRestServer{
		remoteUri: remoteUri,
	}
}

func (hb *HypervisorRestServer) GetUri(action VirtualMachineAction) (string, error) {
	switch action {
	case CREATE:
		return fmt.Sprintf("%s/vm.create", hb.remoteUri), nil
	case BOOT:
		return fmt.Sprintf("%s/vm.boot", hb.remoteUri), nil
	case DELETE:
		return fmt.Sprintf("%s/vm.delete", hb.remoteUri), nil
	case SHUTDOWN:
		return fmt.Sprintf("%s/vm.shutdown", hb.remoteUri), nil
	case INFO:
		return fmt.Sprintf("%s/vm.info", hb.remoteUri), nil
	default:
		return "", errors.New("unknow action")
	}
}
