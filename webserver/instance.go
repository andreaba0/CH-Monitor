package webserver

import (
	"net/http"
	"vmm/vmm"

	"github.com/labstack/echo/v4"
)

func Run(vmmManager *vmm.HypervisorMonitor, socket string) {
	var e *echo.Echo = echo.New()
	var virtualMachineUpload *VirtualMachineUpload = NewVirtualMachineUpload(vmmManager)
	var virtualMachineManagerApi *VirtualMachineManagerApi = NewVirtualMachineManagerApi(vmmManager)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusAccepted, "pong")
	})

	e.POST("/api/disk/upload/:filename/begin", virtualMachineUpload.UploadBegin(UploadType(DISK)))
	e.PUT("/api/disk/upload/:filename/chunk", virtualMachineUpload.UploadChunk(UploadType(DISK)))
	e.POST("/api/disk/upload/:filename/commit", virtualMachineUpload.UploadCommit(UploadType(DISK)))

	e.POST("/api/kernel/upload/:filename/begin", virtualMachineUpload.UploadBegin(UploadType(KERNEL)))
	e.PUT("/api/kernel/upload/:filename/chunk", virtualMachineUpload.UploadChunk(UploadType(KERNEL)))
	e.POST("/api/kernel/upload/:filename/commit", virtualMachineUpload.UploadCommit(UploadType(KERNEL)))

	e.GET("/api/vm/:vm/info", nil)
	e.PUT("/api/vm/:vm/boot", virtualMachineManagerApi.BootVirtualMachine())
	e.PUT("/api/vm/:vm/shutdown", nil)
	e.PUT("/api/vm/:vm/delete", nil)

	e.PUT("/api/vmm/metadata", virtualMachineManagerApi.UpdateVirtualMachine())
	e.POST("/api/vmm/metadata", virtualMachineManagerApi.CreateVirtualMachine())

	e.Logger.Fatal(e.Start(socket))
}
