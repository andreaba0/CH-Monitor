package main

import (
	"flag"
	"log"
	"os"
	virtualmachine "vmm/virtual_machine"
	vmmanager "vmm/vmm"
	"vmm/webserver"

	"go.uber.org/zap"
)

func main() {
	var err error
	var vmFileSystemStorage *virtualmachine.FileSystemWrapper
	var runningCHInstances []vmmanager.RunningCHInstance
	var hostManifestPath string = "/etc/vmm/manifest.json"
	var hostManifest *vmmanager.Manifest
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var hypervisorBinary *vmmanager.HypervisorBinary = nil
	var hypervisorMonitor *vmmanager.HypervisorMonitor

	flag.StringVar(&hostManifestPath, "manifest_path", hostManifestPath, "Path to host manifest")
	flag.Parse()

	hostManifest, err = vmmanager.LoadManifest(hostManifestPath)
	if err != nil {
		logger.Fatal("Unable to load manifest file", zap.String("path", hostManifestPath))
	}

	hypervisorBinary = &vmmanager.HypervisorBinary{
		BinaryPath: hostManifest.HypervisorPath,
		RemoteUri:  hostManifest.HypervisorSocketUri,
	}

	err = os.MkdirAll(hostManifest.Server.StoragePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to create %s folder", hostManifest.Server.StoragePath)
	}

	// Load the list of all active CH processes on the system
	runningCHInstances, err = vmmanager.LoadProcessData(hypervisorBinary)
	if err != nil {
		logger.Fatal("Unable to load processes data")
	}

	// Create an object that is the only one authorized to interact with the filesystem
	vmFileSystemStorage, err = virtualmachine.NewFileSystemWrapper(hostManifest.Server.StoragePath, logger)
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
	webserver.Run(vmFileSystemStorage, hypervisorMonitor, hostManifest.Server.ListeningAddress)

	os.Exit(0)
}
