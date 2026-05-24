package webserver

import (
	"fmt"
	"net/http"
	virtualmachine "vmm/virtual_machine"
	"vmm/vmm"

	"github.com/labstack/echo/v4"
)

type VirtualMachineManagerApi struct {
	vmm *vmm.HypervisorMonitor
}

func NewVirtualMachineManagerApi(vmm *vmm.HypervisorMonitor) *VirtualMachineManagerApi {
	return &VirtualMachineManagerApi{
		vmm: vmm,
	}
}

func (vmmApi *VirtualMachineManagerApi) CreateVirtualMachine() echo.HandlerFunc {
	return func(c echo.Context) error {
		manifest := new(virtualmachine.Manifest)
		var err error
		if err = c.Bind(&manifest); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("There was an error with request body\n%s", err.Error()))
		}
		err = vmmApi.vmm.CreateVirtualMachine(manifest)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("There was an error creating the vm\n%s", err.Error()))
		}
		return c.String(http.StatusAccepted, "Created")
	}
}

func (vmmApi *VirtualMachineManagerApi) BootVirtualMachine() echo.HandlerFunc {
	return func(c echo.Context) error {
		vmId := c.Param("vm")
		virtualMachine := vmmApi.vmm.GetVirtualMachine(vmId)
		if virtualMachine == nil {
			return c.String(http.StatusNotFound, "Virtual Machine is not found")
		}
		err := virtualMachine.RequestBoot(vmmApi.vmm.GetBinaryPath(), vmmApi.vmm.GetRestServerUri())
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("There was a problem booting the vm\n%s", err.Error()))
		}
		return c.String(http.StatusOK, "Booted")
	}
}

func (vmmApi *VirtualMachineManagerApi) UpdateVirtualMachine() echo.HandlerFunc {
	return func(e echo.Context) error {
		return nil
	}
}

type VirtualMachineManagerApiService interface {
	CreateVirtualMachine() echo.HandlerFunc
	UpdateVirtualMachine() echo.HandlerFunc
	BootVirtualMachine() echo.HandlerFunc
}
