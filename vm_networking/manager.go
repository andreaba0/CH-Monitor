package vmnetworking

import (
	"net"

	"github.com/vishvananda/netlink"
)

type VirtualMachineNetworkUtility struct {
	defaultBridgeNetwork net.IPNet
}

type NamingConvention struct {
	Prefix     string
	ObjectName string
	Ip         net.IP
	Mask       net.IPMask
}

type NetworkIdentifier struct {
	Ip        net.IP
	Mask      net.IPMask
	Tenant    string
	GuestName string
}

type NetworkDeviceNameGeneration struct {
	Bridge *string
	Tap    *string
	VxLan  *string
}

func (nc *NamingConvention) Is(value string) bool {
	if nc == nil {
		return false
	}
	return nc.Prefix == value
}

func (nc *NamingConvention) IsVirtualMachineTap() bool {
	return nc.Is("chtap")
}

func NewVirtualMachineNetworkUtility(defaultBridgeNetwork string) (*VirtualMachineNetworkUtility, error) {
	_, ipNet, err := net.ParseCIDR(defaultBridgeNetwork)
	if err != nil {
		return nil, err
	}
	return &VirtualMachineNetworkUtility{
		defaultBridgeNetwork: *ipNet,
	}, nil
}

func filterLinksBy(fn func(netlink.Link) bool) ([]netlink.Link, error) {
	all, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}
	var res []netlink.Link
	for _, link := range all {
		if fn(link) {
			res = append(res, link)
		}
	}
	return res, nil
}

func (nm *VirtualMachineNetworkUtility) GetTapDevices(vmId string) ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		tuntap, ok := link.(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			return false
		}
		return true
	})
}

func (nm *VirtualMachineNetworkUtility) GetAllTapDevices() ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		tuntap, ok := link.(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			return false
		}
		if link.Attrs().MasterIndex == 0 {
			return false
		}
		_, err := netlink.LinkByIndex(link.Attrs().MasterIndex)
		if err != nil {
			return false
		}
		return true
	})
}
