package cloudhypervisor

import (
	"errors"
	"fmt"
	"vmm/utils"
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
	fmt.Printf("remote uri: %s\n", remoteUri)
	return &HypervisorRestServer{
		remoteUri: remoteUri,
	}
}

func (hb *HypervisorRestServer) GetUri(action VirtualMachineAction) (string, error) {
	switch action {
	case CREATE:
		return utils.JoinUri(hb.remoteUri, "/vm.create"), nil
	case BOOT:
		return utils.JoinUri(hb.remoteUri, "/vm.boot"), nil
	case DELETE:
		return utils.JoinUri(hb.remoteUri, "/vm.delete"), nil
	case SHUTDOWN:
		return utils.JoinUri(hb.remoteUri, "/vm.shutdown"), nil
	case INFO:
		return utils.JoinUri(hb.remoteUri, "/vm.info"), nil
	default:
		return "", errors.New("unknow action")
	}
}
