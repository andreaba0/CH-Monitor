package webserver

import (
	"vmm/vmm"

	"github.com/labstack/echo/v4"
)

func Run(vmmManager *vmm.HypervisorMonitor, socket string) {
	var e *echo.Echo = echo.New()
	var virtualMachineUpload *VirtualMachineUpload = NewVirtualMachineUpload(vmmManager)
	var virtualMachineManagerApi *VirtualMachineManagerApi = NewVirtualMachineManagerApi(vmmManager)

	e.PUT("/api/v1/disk/upload/:filename/chunk", virtualMachineUpload.UploadChunk())
	e.POST("/api/v1/disk/upload/:filename/commit", virtualMachineUpload.UploadCommit())

	e.GET("/api/vm/info", nil)
	e.PUT("/api/vm/create", nil)
	e.PUT("/api/vm/boot", nil)
	e.PUT("/api/vm/delete", nil)

	e.PUT("/api/vmm/vm.metadata", virtualMachineManagerApi.UpdateVirtualMachine())
	e.POST("/api/vmm/vm.metadata", virtualMachineManagerApi.CreateVirtualMachine())

	e.Logger.Fatal(e.Start(socket))
}
