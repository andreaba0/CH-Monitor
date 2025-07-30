package vmmanager

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/exec"
	"syscall"
	"time"
	vmnetworking "vmm/vm_networking"

	"go.uber.org/zap"
)

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
	FileSystemStorage *FileSystemWrapper
	NetworkStack      *vmnetworking.VirtualMachineNetworkUtility
	Logger            *zap.Logger
	HypervisorBinary  *HypervisorBinary
	HttpClient        *http.Client
}

func NewVirtualMachine(manifest *Manifest, fileSystemStorage *FileSystemWrapper) (*VirtualMachine, error) {
	var vm *VirtualMachine = &VirtualMachine{
		Manifest: manifest,
	}
	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", fileSystemStorage.GetSocketPath(manifest.GuestName))
	}
	transport := &http.Transport{
		DialContext: dialer,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
	vm.HttpClient = client
	return vm, nil
}

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

func (vm *VirtualMachine) Boot() {
	var _ string = vm.FileSystemStorage.GetSocketPath(vm.Manifest.GuestName)
}

func (vm *VirtualMachine) Shutdown() {

}

func (vm *VirtualMachine) GetState() VirtualMachineState {
	return VirtualMachineState(Unknow)
}
