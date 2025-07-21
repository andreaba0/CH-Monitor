package webserver

import (
	"fmt"
	vmmanager "vmm/manager"

	"github.com/labstack/echo/v4"
)

type EchoSocket struct {
	Port string
}

func Run(vmFileSystemStorage *vmmanager.FileSystemWrapper, vmmManager *vmmanager.HypervisorMonitor, echoSocket *EchoSocket) {
	var e *echo.Echo = echo.New()
	var virtualMachineUpload *VirtualMachineUpload = &VirtualMachineUpload{
		VmFileSystemStorage: vmFileSystemStorage,
	}
	var virtualMachineManagerApi *VirtualMachineManagerApi = &VirtualMachineManagerApi{
		fs:  vmFileSystemStorage,
		vmm: vmmManager,
	}

	e.PUT("/api/v1/disk/upload/:filename/chunk", virtualMachineUpload.UploadChunk())
	e.POST("/api/v1/disk/upload/:filename/commit", virtualMachineUpload.UploadCommit())

	e.GET("/api/vm/info", nil)
	e.PUT("/api/vm/create", nil)
	e.PUT("/api/vm/boot", nil)
	e.PUT("/api/vm/delete", nil)

	e.PUT("/api/vmm/vm.metadata", virtualMachineManagerApi.UpdateVirtualMachine())
	e.POST("/api/vmm/vm.metadata", virtualMachineManagerApi.CreateVirtualMachine())

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", echoSocket.Port)))
}
