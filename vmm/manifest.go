package vmm

import (
	"errors"
	"net"
	vmstorage "vmm/storage"
)

type Manifest struct {
	Bridge              string    `json:"bridge" yaml:"bridge"`
	Server              Server    `json:"server" yaml:"server"`
	Network             net.IPNet `json:"network" yaml:"network"`
	HypervisorPath      string    `json:"hypervisor_path" yaml:"hypervisor_path"`
	HypervisorSocketUri string    `json:"socket_uri" yaml:"socket_uri"`
}

type Server struct {
	StoragePath string `json:"storage_path" yaml:"storage_path"`
}

func LoadManifest(path string) (*Manifest, error) {
	var manifest *Manifest
	var err error

	manifest, err = vmstorage.ReadJson[*Manifest](path)
	if err != nil {
		return nil, errors.New("unable to read manifest file")
	}
	return manifest, nil
}
