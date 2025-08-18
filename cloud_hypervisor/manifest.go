package cloudhypervisor

type Manifest struct {
	Cpus     VmCpus   `json:"cpus" yaml:"cpus"`
	Payload  Payload  `json:"payload" yaml:"payload"`
	Disks    []Disk   `json:"disks" yaml:"disks"`
	Rng      Rng      `json:"rng" yaml:"rng"`
	Net      []Net    `json:"net" yaml:"net"`
	Serial   Serial   `json:"serial" yaml:"serial"`
	Console  Console  `json:"console" yaml:"console"`
	Platform Platform `json:"platform" yaml:"platform"`
}

type Platform struct {
	Uuid string `json:"uuid" yaml:"uuid"`
}

type VmCpus struct {
	Boot_vcpus int `json:"boot_vcpus" yaml:"boot_vcpus"`
	Max_vcpus  int `json:"max_vcpus" yaml:"max_vcpus"`
}

type Payload struct {
	Kernel  string `json:"kernel" yaml:"kernel"`
	Cmdline string `json:"cmdline" yaml:"cmdline"`
}

type Disk struct {
	Path string `json:"path" yaml:"path"`
}

type Rng struct {
	Src string `json:"src" yaml:"src"`
}

type Net struct {
	Tap string `json:"tap" yaml:"tap"`
	Mac string `json:"mac" yaml:"mac"`
}

type Serial struct {
	Mode string `json:"mode" yaml:"mode"`
	File string `json:"file" yaml:"file"`
}

type Console struct {
	Mode string `json:"mode" yaml:"mode"`
}
