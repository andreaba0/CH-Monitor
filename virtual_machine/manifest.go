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

func (manifest *Manifest) ParseToInstanceRequest() (*cloudhypervisor.CloudHypervisor, error) {
	return nil, nil
}

type Config struct {
	Networks []Net  `json:"networks" yaml:"networks"`
	Disks    []Disk `json:"disks" yaml:"disks"`
}

type Net struct {
	Address string `json:"address" yaml:"address"`
	Mac     string `json:"mac" yaml:"mac"`
}

func (n *Net) ParseToInstanceRequest() (*cloudhypervisor.Net, error) {
	if n.Address == "" || n.Mac == "" {
		return nil, errors.New("properity not filled")
	}
	ip, ipNet, err := net.ParseCIDR(n.Address)
	if err != nil {
		return nil, errors.New("there was an error parsing ip address")
	}
	var tap string = vmnetworking.GenerateTapName(ip, *ipNet)
	var mask string = net.IP(ipNet.Mask).To4().String()
	return &cloudhypervisor.Net{
		Ip:   ip.To4().String(),
		Mask: mask,
		Mac:  n.Mac,
		Tap:  tap,
	}, nil
}

type Disk struct {
	Name string `json:"name" yaml:"name"`
}
