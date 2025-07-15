package webserver

import (
	"fmt"
	vmstorage "vmm/vm_storage"

	"github.com/labstack/echo/v4"
)

type EchoSocket struct {
	Port string
}

func Run(vmFileSystemStorage *vmstorage.FileSystemStorage, echoSocket *EchoSocket) {
	var e *echo.Echo = echo.New()
	var virtualMachineUpload VirtualMachineUploadService = &VirtualMachineUpload{
		VmFileSystemStorage: vmFileSystemStorage,
	}

	e.POST("/api/v1/disk/upload/:filename/begin", virtualMachineUpload.UploadBegin())
	e.PUT("/api/v1/disk/upload/:filename/chunk", virtualMachineUpload.UploadChunk())
	e.POST("/api/v1/disk/upload/:filename/commit", virtualMachineUpload.UploadCommit())

	e.GET("/api/vm/info", nil)
	e.PUT("/api/vm/create", nil)
	e.PUT("/api/vm/boot", nil)
	e.PUT("/api/vm/delete", nil)

	e.PUT("/api/vmm/config/reload", nil)
	e.POST("/api/vmm/vm.metadata", nil)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", echoSocket.Port)))
}
