package webserver

import (
	"vmm/vmm"

	"github.com/labstack/echo/v4"
)

func Run(vmmManager *vmm.HypervisorMonitor, socket string) {
	var e *echo.Echo = echo.New()
	var virtualMachineUpload *VirtualMachineUpload = NewVirtualMachineUpload(vmmManager)
	var virtualMachineManagerApi *VirtualMachineManagerApi = NewVirtualMachineManagerApi(vmmManager)

	e.POST("/api/disk/upload/:filename/begin", virtualMachineUpload.UploadBegin(UploadType(DISK)))
	e.PUT("/api/disk/upload/:filename/chunk", virtualMachineUpload.UploadChunk())
	e.POST("/api/disk/upload/:filename/commit", virtualMachineUpload.UploadCommit())

	e.POST("/api/kernel/upload/:kernelname/begin", virtualMachineUpload.UploadBegin(UploadType(KERNEL)))
	e.PUT("/api/kernel/upload/:kernelname/chunk", nil)
	e.POST("/api/kernel/upload/:kernelname/commit", nil)

	e.GET("/api/vm/info", nil)
	e.PUT("/api/vm/create", nil)
	e.PUT("/api/vm/boot", nil)
	e.PUT("/api/vm/shutdown", nil)
	e.PUT("/api/vm/delete", nil)

	e.PUT("/api/vmm/metadata", virtualMachineManagerApi.UpdateVirtualMachine())
	e.POST("/api/vmm/metadata", virtualMachineManagerApi.CreateVirtualMachine())

	e.Logger.Fatal(e.Start(socket))
}
