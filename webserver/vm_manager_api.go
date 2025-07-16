package webserver

import (
	vmmanager "vmm/vm_manager"
	vmstorage "vmm/vm_storage"

	"github.com/labstack/echo/v4"
)

type VirtualMachineManagerApi struct {
	fs  vmstorage.FileSystemStorageService
	vmm vmmanager.HypervisorMonitorService
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
