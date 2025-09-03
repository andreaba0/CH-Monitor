package vmm

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	Bridge                   string `json:"bridge" yaml:"bridge"`
	Server                   Server `json:"server" yaml:"server"`
	HypervisorPath           string `json:"hypervisor_path" yaml:"hypervisor_path"`
	HypervisorSocketUri      string `json:"socket_uri" yaml:"socket_uri"`
	InternalConfigFolderPath string `json:"config_folder_path" yaml:"config_folder_path"`
}

type Server struct {
	StoragePath string `json:"storage_path" yaml:"storage_path"`
}

func LoadManifest(path string) (*Manifest, error) {
	manifest := &Manifest{}

	fileByte, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileByte, manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}
