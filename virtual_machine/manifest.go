package virtualmachine

import (
	"github.com/google/uuid"
)

type Manifest struct {
	GuestIdentifier uuid.UUID `json:"guest_identifier" xml:"guest_identifier"`
	Tenant          uuid.UUID `json:"tenant" xml:"tenant"`
	Config          Config    `json:"hypervisor_config" yaml:"hypervisor_config"`
}

type Config struct {
	Network Net    `json:"networks" yaml:"networks"`
	Disks   []Disk `json:"disks" yaml:"disks"`
	Kernel  string `json:"kernel" yaml:"kernel"`
	Init    string `json:"init" yaml:"init"`
	Vpc     []Net  `json:"vpc" yaml:"vpc"`
	Rng     Rng    `json:"rng" yaml:"rng"`
	Cpus    int    `json:"cpus" yaml:"cpus"`
}

type Rng struct {
	Src string `json:"src" yaml:"src"`
}

type Net struct {
	Addresses []string `json:"addresses" yaml:"addresses"`
	Mask      string   `json:"mask" yaml:"mask"`
	Mac       string   `json:"mac" yaml:"mac"`
	Tap       string   `json:"tap" yaml:"tap"`
}

type Disk struct {
	Name string `json:"name" yaml:"name"`
}
