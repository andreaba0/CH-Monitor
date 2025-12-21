package main

import (
	"flag"
	"os"
	vmmanager "vmm/vmm"
	"vmm/webserver"

	"go.uber.org/zap"
)

func main() {
	var err error
	var hostManifestPath string
	var serverAddress string
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	flag.StringVar(&hostManifestPath, "manifest_path", "/etc/vmm/manifest.json", "Path to host manifest")
	flag.StringVar(&serverAddress, "server_address", "0.0.0.0:8080", "Address to bind")
	flag.Parse()

	hypervisorMonitor, err := vmmanager.NewHypervisorMonitor(logger, hostManifestPath)
	if err != nil {
		logger.Fatal("unable to initialize vmm", zap.String("error", err.Error()))
	}

	err = hypervisorMonitor.MonitorSetup(hostManifestPath, hypervisorMonitor)
	if err != nil {
		logger.Fatal("Unable to init vmm", zap.String("error", err.Error()))
	}

	// Run webserver and start listening for incoming requests
	webserver.Run(hypervisorMonitor, serverAddress)

	os.Exit(0)
}
