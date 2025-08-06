package webserver

import (
	"fmt"
	"net/http"
	virtualmachine "vmm/virtual_machine"
	"vmm/vmm"

	"github.com/labstack/echo/v4"
)

type UploadType int

const (
	KERNEL UploadType = iota
	DISK
)

type BeginBody struct {
	VirtualMachine string `json:"virtual_machine" xml:"virtual_machine"`
}

type CommitBody struct {
	VirtualMachine string `json:"virtual_machine" xml:"virtual_machine"`
	TmpDiskName    string `json:"tmp_disk_name" xml:"tmp_disk_name"`
}

type JsonResponse struct {
	Message string `json:"message" xml:"message"`
}

type VirtualMachineUpload struct {
	vmm *vmm.HypervisorMonitor
}

func NewVirtualMachineUpload(vmm *vmm.HypervisorMonitor) *VirtualMachineUpload {
	return &VirtualMachineUpload{
		vmm: vmm,
	}
}

type VirtualMachineUploadService interface {
	UploadBegin() echo.HandlerFunc
	UploadCommit() echo.HandlerFunc
	UploadChunk() echo.HandlerFunc
	CreateVirtualMachine() echo.HandlerFunc
}

func (vmStorage *VirtualMachineUpload) UploadBegin(uploadType UploadType) echo.HandlerFunc {
	return func(c echo.Context) error {
		fileMetadata := new(BeginBody)
		var err error
		if err = c.Bind(fileMetadata); err != nil {
			return c.String(http.StatusBadRequest, "Malformed request body")
		}
		var filename = c.Param("filename")
		var vm *virtualmachine.VirtualMachine = vmStorage.vmm.GetVirtualMachine(fileMetadata.VirtualMachine)
		if vm == nil {
			return c.String(http.StatusNotFound, "Requested virtual machine is not found")
		}

		var tmpFileName string
		if uploadType == UploadType(DISK) {
			tmpFileName, err = vm.CreateDisk(filename)
			if err != nil {
				return c.String(http.StatusBadRequest, "There was an error creating disk")
			}
		} else if uploadType == UploadType(KERNEL) {
			tmpFileName, err = vm.CreateKernel(filename)
			if err != nil {
				return c.String(http.StatusBadRequest, "There was an error creating kernel")
			}
		} else {
			return c.String(http.StatusBadRequest, "Unknow file kind")
		}
		return c.JSON(http.StatusOK, JsonResponse{Message: tmpFileName})
	}
}

func (vmStorage *VirtualMachineUpload) UploadCommit() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileMetadata := new(CommitBody)
		var err error
		if err = c.Bind(fileMetadata); err != nil {
			return c.String(http.StatusBadRequest, "Provided request body does not fulfill requirements")
		}
		var filename = c.Param("filename")

		var vm *virtualmachine.VirtualMachine = vmStorage.vmm.GetVirtualMachine(fileMetadata.VirtualMachine)
		if vm == nil {
			return c.String(http.StatusNotFound, "Virtual Machine is not found")
		}

		err = vm.CommitDisk(fileMetadata.TmpDiskName, filename)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, JsonResponse{Message: "There was an error completing file upload"})
		}
		return c.JSON(http.StatusOK, JsonResponse{Message: "OK"})
	}
}

func (vmStorage *VirtualMachineUpload) UploadChunk() echo.HandlerFunc {
	return func(c echo.Context) error {
		var header http.Header
		var filename string = c.Param("filename")
		var contentRange string
		var rangeStart int64
		var rangeEnd int64
		var fileSize int64
		var sizeUnit string
		var err error
		var n int
		var virtualMachine = c.Request().Header.Get("X-VirtualMachine")
		header = c.Request().Header
		contentRange = header.Get("Content-Range")
		if contentRange == "" {
			return c.JSON(http.StatusBadRequest, JsonResponse{Message: "Content-Range required for chunked upload"})
		}
		n, err = fmt.Sscanf(contentRange, "%s %d-%d/%d", &sizeUnit, &rangeStart, &rangeEnd, &fileSize)
		if err != nil || n <= 0 {
			return c.JSON(http.StatusBadRequest, JsonResponse{Message: "Maybe Content-Range is in bad format. Required: <unit> <range start>-<range end>/<full size>"})
		}

		var vm *virtualmachine.VirtualMachine = vmStorage.vmm.GetVirtualMachine(virtualMachine)
		if vm == nil {
			return c.String(http.StatusNotFound, "Virtual Machine is not found")
		}

		err = vm.WriteChunkToDisk(filename, rangeStart, c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, JsonResponse{Message: "There was an error writing to file"})
		}

		return c.JSON(http.StatusOK, JsonResponse{Message: "OK"})
	}
}
