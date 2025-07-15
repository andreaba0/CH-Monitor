package webserver

import (
	"fmt"
	"net/http"
	vmstorage "vmm/vm_storage"

	"github.com/labstack/echo/v4"
)

type DiskMetadata struct {
	VirtualMachine string `json:"virtual_machine" xml:"virtual_machine"`
	ByteSize       int64  `json:"byte_size" xml:"byte_size"`
}

type JsonResponse struct {
	Message string `json:"message" xml:"message"`
}

type VirtualMachineUpload struct {
	VmFileSystemStorage *vmstorage.FileSystemStorage
}

type VirtualMachineUploadService interface {
	UploadBegin() echo.HandlerFunc
	UploadCommit() echo.HandlerFunc
	UploadChunk() echo.HandlerFunc
	CreateVirtualMachine() echo.HandlerFunc
}

func (vmStorage *VirtualMachineUpload) CreateVirtualMachine() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (vmStorage *VirtualMachineUpload) UploadBegin() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileMetadata := new(DiskMetadata)
		var vmName string = c.Param("filename")
		var err error
		if err = c.Bind(fileMetadata); err != nil {
			return c.String(http.StatusBadRequest, "Provided request body does not fulfill requirements")
		}
		err = vmStorage.VmFileSystemStorage.CreateDisk(vmName, fileMetadata.VirtualMachine, fileMetadata.ByteSize)
		if err != nil {
			return c.JSON(http.StatusInsufficientStorage, JsonResponse{Message: "Unable to allocate enough space for file"})
		}

		return c.JSON(http.StatusOK, JsonResponse{Message: "BEGIN UPLOAD"})
	}
}

func (vmStorage *VirtualMachineUpload) UploadCommit() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileMetadata := new(DiskMetadata)
		var err error
		if err = c.Bind(fileMetadata); err != nil {
			return c.String(http.StatusBadRequest, "Provided request body does not fulfill requirements")
		}
		var filename = c.Param("filename")
		err = vmStorage.VmFileSystemStorage.CompleteDiskWrite(fileMetadata.VirtualMachine, filename)
		if err != nil {
			return c.JSON(http.StatusBadRequest, JsonResponse{Message: "Maybe file does not exists"})
		}
		return c.String(http.StatusOK, "COMMIT")
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

		err = vmStorage.VmFileSystemStorage.WriteDiskChunk(virtualMachine, filename, rangeStart, c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, JsonResponse{Message: "There was an error writing to file"})
		}

		return c.JSON(http.StatusOK, JsonResponse{Message: "OK"})
	}
}
