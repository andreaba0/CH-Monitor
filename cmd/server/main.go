package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	vmmanager "vmm/manager"
	"vmm/webserver"

	"go.uber.org/zap"
)

func main() {
	var err error
	var homeDir string
	var basePath = ""
	var vmFileSystemStorage *vmmanager.FileSystemWrapper
	var runningCHInstances []vmmanager.RunningCHInstance
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var hypervisorBinary *vmmanager.HypervisorBinary = &vmmanager.HypervisorBinary{
		BinaryPath: "/bin/cloud-hypervisor-static",
		RemoteUri:  "http://localhost",
	}
	var defaultNetwork string = "192.168.0.0/24"
	var hypervisorMonitor *vmmanager.HypervisorMonitor
	var echoSocket *webserver.EchoSocket = &webserver.EchoSocket{}

	homeDir, err = os.UserHomeDir()
	if err != nil {
		basePath = filepath.Join(homeDir, "hypervisor/storage")
	} else {
		basePath = filepath.Join("hypervisor/storage")
	}

	flag.StringVar(&basePath, "storage_path", basePath, "Path to folder where vms disk rootfs are stored")
	flag.StringVar(&echoSocket.Port, "port", "80", "Webserver listening port")
	flag.StringVar(&defaultNetwork, "default_network", defaultNetwork, "Default network to attach guest vms to")

	flag.Parse()

	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to create %s folder", basePath)
	}

	// Load the list of all active CH processes on the system
	runningCHInstances, err = vmmanager.LoadProcessData(hypervisorBinary)
	if err != nil {
		logger.Fatal("Unable to load processes data")
	}

	// Create an object that is the only one authorized to interact with the filesystem
	vmFileSystemStorage, err = vmmanager.NewFileSystemWrapper(basePath, logger)
	if err != nil {
		log.Fatal("unable to load file system storage")
	}

	// Initialize the VMM
	hypervisorMonitor = vmmanager.NewHypervisorMonitor(vmFileSystemStorage, logger)

	vmList, err := vmFileSystemStorage.GetVirtualMachineList()
	if err != nil {
		log.Fatal("Unable to query vm list")
	}
	hypervisorMonitor.LoadVirtualMachines(runningCHInstances, vmList)

	// Run webserver and start listening for incoming requests
	webserver.Run(vmFileSystemStorage, hypervisorMonitor, echoSocket)

	os.Exit(0)
}
