package webserver

import (
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
		if c.Bind(manifest); err != nil {
			return c.String(http.StatusBadRequest, "There was an error with request body")
		}
		err = vmmApi.vmm.CreateVirtualMachine(manifest)
		if err != nil {
			return c.String(http.StatusBadRequest, "There was an error creating the vm")
		}
		return c.String(http.StatusAccepted, "Created")
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
}
