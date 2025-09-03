package vmm

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"
	virtualmachine "vmm/virtual_machine"
	vmnetworking "vmm/vm_networking"
	networkvpc "vmm/vm_networking/vpc"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type HypervisorMonitor struct {
	virtualMachines   map[string]*virtualmachine.VirtualMachine
	vmsMu             sync.Mutex
	logger            *zap.Logger
	manifest          *Manifest
	networkEnumerator *vmnetworking.NetworkEnumerator
	vpcManager        *networkvpc.VpcManager
}

func NewHypervisorMonitor(logger *zap.Logger, manifestPath string) (*HypervisorMonitor, error) {
	fileBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	manifest := &Manifest{}
	err = yaml.Unmarshal(fileBytes, manifest)
	if err != nil {
		return nil, err
	}
	enumeratorFilePath := filepath.Join(manifest.InternalConfigFolderPath, "enumerator_config.json")
	vpcSnapshotFilePath := filepath.Join(manifest.InternalConfigFolderPath, "vpc_config.snapshot")
	vpcChangesFilePath := filepath.Join(manifest.InternalConfigFolderPath, "vpc_changes.aof")
	networkEnumerator, err := vmnetworking.NewNetworkEnumerator(enumeratorFilePath)
	if err != nil {
		return nil, err
	}
	return &HypervisorMonitor{
		virtualMachines:   make(map[string]*virtualmachine.VirtualMachine),
		logger:            logger,
		manifest:          manifest,
		networkEnumerator: networkEnumerator,
		vpcManager:        networkvpc.NewVpcManager(vpcSnapshotFilePath, vpcChangesFilePath),
	}, nil
}

func (hm *HypervisorMonitor) MonitorSetup(manifestPath string, vmm *HypervisorMonitor) error {
	var hypervisorBinary cloudhypervisor.HypervisorRestServer = *cloudhypervisor.NewHypervisorRestServer(hm.manifest.HypervisorSocketUri)
	err := vmm.LoadVirtualMachines(hm.manifest.Server.StoragePath)
	if err != nil {
		return err
	}
	err = vmm.MergeRunningInstances(hm.manifest.HypervisorPath, &hypervisorBinary)
	if err != nil {
		return err
	}
	return nil
}

func (hm *HypervisorMonitor) SetManifest(manifest *Manifest) {
	hm.manifest = manifest
}

func (hm *HypervisorMonitor) LoadVirtualMachines(basePath string) error {

	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	var err error
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		vm, err := virtualmachine.LoadVirtualMachine(filepath.Join(basePath, entry.Name()), hm.logger, hm.manifest.Bridge, hm.networkEnumerator)
		if err != nil {
			hm.logger.Error("Unable to read manifest from file", zap.String("base_path", basePath), zap.String("vm_id", entry.Name()))
		}
		guestName := vm.GetManifest().GuestIdentifier
		hm.virtualMachines[guestName.String()] = vm
	}
	return nil
}

func (hm *HypervisorMonitor) MergeRunningInstances(hypervisorBinaryPath string, hypervisorBinary *cloudhypervisor.HypervisorRestServer) error {
	var err error
	var instances []*cloudhypervisor.CloudHypervisor
	instances, err = LoadProcessData(hypervisorBinaryPath)
	if err != nil {
		return err
	}
	for i := 0; i < len(instances); i++ {
		var manifest cloudhypervisor.Manifest
		var err error
		uri, err := hypervisorBinary.GetUri(cloudhypervisor.VirtualMachineAction(cloudhypervisor.INFO))
		if err != nil {
			return err
		}
		res, err := instances[i].HttpClient.Get(uri)
		if err != nil {
			return err
		}
		if res.StatusCode < 200 || res.StatusCode > 299 {
			return errors.New("error while retrieving vm info")
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(resBody, &manifest)
		if err != nil {
			return err
		}
		var vm *virtualmachine.VirtualMachine
		hm.vmsMu.Lock()
		vm = hm.virtualMachines[manifest.Platform.Uuid]
		hm.vmsMu.Unlock()
		vm.AttachInstance(instances[i])
	}
	return nil
}

func (hm *HypervisorMonitor) CreateVirtualMachine(manifest *virtualmachine.Manifest) error {
	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	for i := 0; i < len(manifest.Config.Vpc); i++ {
		vpc := manifest.Config.Vpc[i]
		if len(vpc.Addresses) < 1 {
			return errors.New("required at least one ip address for a given interface")
		}
		if len(vpc.Addresses) < 1 {
			return errors.New("expected at least 1 ip address")
		}
		_, ipNet, err := vmnetworking.ParseCIDR4(vpc.Addresses[0], vpc.Mask)
		if err != nil {
			return err
		}
		bridge, err := hm.networkEnumerator.GetNewBridgeName()
		if err != nil {
			return err
		}
		err = hm.vpcManager.AddNetwork(manifest.Tenant, *ipNet, bridge)
		if err != nil {
			return err
		}
		manifest.Config.Vpc[i].Bridge = bridge
	}
	return nil
}

func (hm *HypervisorMonitor) GetVirtualMachine(id string) *virtualmachine.VirtualMachine {
	hm.vmsMu.Lock()
	defer hm.vmsMu.Unlock()
	return hm.virtualMachines[id]
}

func (hm *HypervisorMonitor) GetBinaryPath() string {
	return hm.manifest.HypervisorPath
}

func (hm *HypervisorMonitor) GetRestServerUri() string {
	return hm.manifest.HypervisorSocketUri
}
