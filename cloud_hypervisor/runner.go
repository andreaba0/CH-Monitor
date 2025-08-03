package cloudhypervisor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/vishvananda/netlink"
)

type CloudHypervisor struct {
	pid        int
	httpClient *http.Client
}

func CreateTransportSocket(socket string) *http.Client {
	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", socket)
	}
	transport := &http.Transport{
		DialContext: dialer,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
	return client
}

func NewCloudHypervisor(manifest *Manifest, binaryPath string, defaultBridge netlink.Link) (*CloudHypervisor, error) {

	// TODO:
	// 1. Create tap interfaces based on manifest networks and attach them to default bridge
	// 2. Launch cloud-hypervisor instance as daemon
	// 3. Return launched instance

	socketUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.New("error while generating uuid for socket file")
	}
	var socketPath string = fmt.Sprintf("/tmp/vm-net-%s.sock", socketUuid)
	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		httpClient: CreateTransportSocket(socketPath),
	}
	return cloudHypervisor, nil
}

func (ch *CloudHypervisor) Kill(manifest *Manifest) error {

	// TODO:
	// 1. Terminate current cloud-hypervisor instance
	// 2. Drop all tap interfaces connected to the instance

	proc, err := os.FindProcess(ch.pid)
	if err != nil {
		return errors.New("there was an error searching process by pid")
	}
	err = proc.Signal(syscall.SIGKILL)
	if err != nil {
		return errors.New("there was an error killing running process")
	}
	return nil
}

func LoadRunningInstance(pid int, socketPath string) *CloudHypervisor {
	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		pid:        pid,
		httpClient: CreateTransportSocket(socketPath),
	}
	return cloudHypervisor
}

/*
func (vm *VirtualMachine) Create() error {
	master, err := vm.NetworkStack.GetDefaultBridge()
	if err != nil {
		return err
	}
	for i := 0; i < len(vm.Manifest.Net); i++ {
		currentNet := vm.Manifest.Net[i]
		interfaceName, err := vm.NetworkStack.GenereateDeviceName(&vmnetworking.NetworkIdentifier{
			Ip:        currentNet.Address.IP,
			Mask:      currentNet.Address.IPNet.Mask,
			Tenant:    vm.Manifest.Tenant,
			GuestName: vm.Manifest.GuestName,
		})
		if err != nil {
			vm.Logger.Error("There was an error creating network interface name")
			continue
		}
		err = vmnetworking.CreateTapInterface(*interfaceName.Tap, currentNet.Address.IP, currentNet.Address.IPNet.Mask, master)
		if err != nil {
			vm.Logger.Error("There was an error creating a tap interface",
				zap.String("guest_vm", vm.Manifest.GuestName),
				zap.String("ip", string(currentNet.Address.IP)),
			)
		}
	}

	cmd := exec.Command(vm.HypervisorBinary.GetLaunchCommandString(vm.FileSystemStorage.GetSocketPath(vm.Manifest.GuestName)))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		vm.Logger.Error("Unable to run process", zap.String("guest_name", vm.Manifest.GuestName))
		return err
	}
	vm.PID = &cmd.Process.Pid

	var uri *string = vm.HypervisorBinary.GetUri(CREATE)
	if uri == nil {
		vm.Logger.Error("Unknow URI", zap.String("action", "CREATE"))
		return errors.New("unknow uri")
	}
	req, err := http.NewRequest("POST", *uri, nil)
	if err != nil {
		return err
	}

	resp, err := vm.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
*/
