package vmstorage

import (
	"net"
)

type IPNetWrapper struct {
	net.IPNet
}

func (ipnet *IPNetWrapper) UnmarshalText(text []byte) error {
	var ip net.IP
	var ipNet *net.IPNet
	var err error
	ip, ipNet, err = net.ParseCIDR(string(text))
	if err != nil {
		return err
	}
	ipnet.IP = ip
	ipnet.Mask = ipNet.Mask
	return nil
}

type NetworkInterfaceConfig struct {
	Addresses []IPNetWrapper   `json:"addresses" xml:"addresses"`
	Mac       net.HardwareAddr `json:"mac" xml:"mac"`
}

type DiskDevice struct {
	Name string `json:"name" xml:"name"`
	Size int64  `json:"size" xml:"size"`
}

type MemoryConfig struct {
	Size int `json:"size" xml:"size"`
}

type CpuConfig struct {
	Count int    `json:"vcpu" xml:"vcpu"`
	Limit string `json:"limit" xml:"limit"`
}

type Manifest struct {
	GuestName         string                   `json:"guest_name" xml:"guest_name"`
	Disks             []DiskDevice             `json:"disks" xml:"disks"`
	EnableIpSpoofing  bool                     `json:"enable_ip_spoofing" xml:"enable_ip_spoofing"`
	EnableMacSpoofing bool                     `json:"enable_mac_spoofing" xml:"enable_mac_spoofing"`
	EnableBroadcast   bool                     `json:"enable_broadcast" xml:"enable_broadcast"`
	Memory            MemoryConfig             `json:"memory" xml:"memory"`
	Vcpu              CpuConfig                `json:"vcpu" xml:"vcpu"`
	Tenant            string                   `json:"tenant" xml:"tenant"`
	DefaultNetwork    NetworkInterfaceConfig   `json:"default_network" xml:"default_network"`
	PrivateNetworks   []NetworkInterfaceConfig `json:"private_networks" xml:"private_networks"` // VPC
}
