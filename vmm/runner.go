package vmm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"sync"
	cloudhypervisor "vmm/cloud_hypervisor"
	"vmm/metadata"
	vmstorage "vmm/storage"
	virtualmachine "vmm/virtual_machine"
	vmnetworking "vmm/vm_networking"

	"go.uber.org/zap"
)

type HypervisorMonitor struct {
	virtualMachines map[string]*virtualmachine.VirtualMachine
	vmsMu           sync.Mutex
	logger          *zap.Logger
	manifest        *Manifest
	metadata        *metadata.MetadataManager
	vpc             map[string]map[string]string // vpc[ip/mask][tenant]=bridge_name
}

func NewHypervisorMonitor(logger *zap.Logger, manifestPath string) (*HypervisorMonitor, error) {
	manifest, err := vmstorage.ReadYaml[Manifest](manifestPath)
	if err != nil {
		return nil, err
	}
	metadataManager, err := metadata.NewMetadataManager(manifest.InternalMetadataPath)
	if err != nil {
		return nil, err
	}
	return &HypervisorMonitor{
		virtualMachines: make(map[string]*virtualmachine.VirtualMachine),
		logger:          logger,
		manifest:        &manifest,
		metadata:        metadataManager,
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
	entries, err := vmstorage.ListFolder(basePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsFolder {
			continue
		}
		vm, err := virtualmachine.LoadVirtualMachine(filepath.Join(basePath, entry.Name), hm.logger, hm.manifest.Bridge, hm.metadata)
		if err != nil {
			hm.logger.Error("Unable to read manifest from file", zap.String("base_path", basePath), zap.String("vm_id", entry.Name))
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
		network := manifest.Config.Vpc[i]
		if len(network.Addresses) < 1 {
			return errors.New("required at least one ip address for a given interface")
		}
		ipNet4 := ""
		ones := 0
		for j := 1; j < len(network.Addresses); j++ {
			_, ipNet, err := vmnetworking.ParseCIDR4(network.Addresses[j], network.Mask)
			ones, _ = ipNet.Mask.Size()
			if err != nil {
				return err
			}
			if ipNet4 == "" {
				ipNet4 = ipNet.String()
				continue
			}
			if ipNet.String() != ipNet4 {
				return errors.New("all ip addresses in an interface must be in the same network")
			}
		}
		tapName, err := hm.metadata.GetNewTapName()
		if err != nil {
			return err
		}
		manifest.Config.Vpc[i].Tap = tapName
		vpc := fmt.Sprintf("%s/%s", ipNet4, strconv.Itoa(ones))
		if v, ok := hm.vpc[vpc][manifest.Tenant.String()]; ok {
			manifest.Config.Vpc[i].Bridge = v
			continue
		}
		bridge, err := hm.metadata.GetNewBridgeName()
		if err != nil {
			return err
		}
		hm.vpc[vpc][manifest.Tenant.String()] = bridge
		manifest.Config.Vpc[i].Bridge = bridge
	}
	tapName, err := hm.metadata.GetNewTapName()
	if err != nil {
		return err
	}
	manifest.Config.Network.Tap = tapName
	vm, err := virtualmachine.NewVirtualMachine(manifest, hm.logger, hm.manifest.Server.StoragePath, hm.manifest.Bridge, hm.metadata)
	if err != nil {
		return err
	}
	err = vm.StoreManifest()
	if err != nil {
		return err
	}
	hm.virtualMachines[manifest.GuestIdentifier.String()] = vm
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
