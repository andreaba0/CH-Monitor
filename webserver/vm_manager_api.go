package webserver

import (
	vmmanager "vmm/manager"

	"github.com/labstack/echo/v4"
)

type VirtualMachineManagerApi struct {
	fs  *vmmanager.FileSystemWrapper
	vmm *vmmanager.HypervisorMonitor
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
