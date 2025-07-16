package vmstorage

type NetworkInterfaceConfig struct {
	Ip      string `json:"ip" xml:"ip"`
	NetMask string `json:"net_mask" xml:"net_mask"`
	Mac     string `json:"mac" xml:"mac"`
}

type DiskDevice struct {
	Name string `json:"name" xml:"name"`
	Size int64  `json:"size" xml:"size"`
}

type MemoryConfig struct {
	Size int `json:"size" xml:"size"`
}

type CpuConfig struct {
	VcpuNumber int    `json:"vcpu" xml:"vcpu"`
	Limit      string `json:"limit" xml:"limit"` // percentage for each core
}

type Manifest struct {
	GuestName         string                   `json:"guest_name" xml:"guest_name"`
	Networks          []NetworkInterfaceConfig `json:"networks" xml:"networks"`
	Disks             []DiskDevice             `json:"disks" xml:"disks"`
	EnableIpSpoofing  bool                     `json:"enable_ip_spoofing" xml:"enable_ip_spoofing"`
	EnableMacSpoofing bool                     `json:"enable_mac_spoofing" xml:"enable_mac_spoofing"`
	EnableBroadcast   bool                     `json:"enable_broadcast" xml:"enable_broadcast"`
	Memory            MemoryConfig             `json:"memory" xml:"memory"`
	Cpu               CpuConfig                `json:"cpu" xml:"cpu"`
}
