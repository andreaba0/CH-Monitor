package cloudhypervisor

import "github.com/google/uuid"

type Manifest struct {
	Cpus     VmCpus   `json:"cpus" yaml:"cpus"`
	Payload  Payload  `json:"payload" yaml:"payload"`
	Disks    []Disks  `json:"disks" yaml:"disks"`
	Rng      Rng      `json:"rng" yaml:"rng"`
	Net      []Net    `json:"net" yaml:"net"`
	Serial   Serial   `json:"serial" yaml:"serial"`
	Console  Console  `json:"console" yaml:"console"`
	Platform Platform `json:"platform" yaml:"platform"`
}

type Platform struct {
	Uuid uuid.UUID `json:"uuid" yaml:"uuid"`
}

type VmCpus struct {
	Boot_vcpus int `json:"boot_vcpus" yaml:"boot_vcpus"`
	Max_vcpus  int `json:"max_vpcus" yaml:"max_vcpus"`
}

type Payload struct {
	Kernel  string `json:"kernel" yaml:"kernel"`
	Cmdline string `json:"cmdline" yaml:"cmdline"`
}

type Disks struct {
	Name string `json:"name" yaml:"name"`
}

type Rng struct {
	Src string `json:"src" yaml:"src"`
}

type Net struct {
	Mac  string `json:"mac" yaml:"mac"`
	Ip   string `json:"ip" yaml:"ip"`
	Mask string `json:"mask" yaml:"mask"`
	Tap  string `json:"tap" yaml:"tap"`
}

type Serial struct {
	Mode string `json:"mode" yaml:"mode"`
}

type Console struct {
	Mode string `json:"mode" yaml:"mode"`
}

// This should return a struct with the same fields as "net" field in cloud-hypervisor vm configuration
func (n *Net) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

func (n *Net) UnmarshalJSON(data []byte) error {
	return nil
}

func (d *Disks) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}
