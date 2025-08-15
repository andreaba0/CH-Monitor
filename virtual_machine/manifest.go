package virtualmachine

import (
	"errors"
	"net"
	cloudhypervisor "vmm/cloud_hypervisor"
	vmnetworking "vmm/vm_networking"

	"github.com/google/uuid"
)

type Manifest struct {
	GuestIdentifier uuid.UUID `json:"guest_identifier" xml:"guest_identifier"`
	Tenant          uuid.UUID `json:"tenant" xml:"tenant"`
	Config          Config    `json:"hypervisor_config" yaml:"hypervisor_config"`
}

type Config struct {
	Networks []Net  `json:"networks" yaml:"networks"`
	Disks    []Disk `json:"disks" yaml:"disks"`
	Kernel   string `json:"kernel" yaml:"kernel"`
	Init     string `json:"init" yaml:"init"`
	Vpc      []Net  `json:"vpc" yaml:"vpc"`
	Rng      Rng    `json:"rng" yaml:"rng"`
	Cpus     int    `json:"cpus" yaml:"cpus"`
}

type Rng struct {
	Src string `json:"src" yaml:"src"`
}

type Net struct {
	Address string `json:"address" yaml:"address"`
	Mac     string `json:"mac" yaml:"mac"`
}

func (n *Net) ParseToInstanceRequest() (*cloudhypervisor.Net, error) {
	if n.Address == "" || n.Mac == "" {
		return nil, errors.New("property not filled")
	}
	ip, ipNet, err := net.ParseCIDR(n.Address)
	if err != nil {
		return nil, errors.New("there was an error parsing ip address")
	}
	var tap string = vmnetworking.GenerateTapName(ip, *ipNet)
	return &cloudhypervisor.Net{
		Mac: n.Mac,
		Tap: tap,
	}, nil
}

type Disk struct {
	Name string `json:"name" yaml:"name"`
}
