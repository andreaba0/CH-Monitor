package webserver

import (
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
		return nil
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
