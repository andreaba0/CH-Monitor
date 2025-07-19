package vmnetworking

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

type VirtualMachineNetworkUtility struct {
	bridgeDefaultName string
}

type NamingConvention struct {
	Prefix     string
	ObjectName string
	Ip         net.IP
	Mask       net.IPMask
}

func ParseDeviceName(name string) (*NamingConvention, error) {
	var parts []string = strings.Split(name, "-")
	if len(parts) < 7 {
		return nil, errors.New("wrong naming convention: format")
	}
	var strMask = parts[len(parts)-1]
	var strIp string = strings.Join([]string{parts[len(parts)-5], parts[len(parts)-4], parts[len(parts)-3], parts[len(parts)-2]}, ".")
	ip, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%s", strIp, strMask))
	if err != nil {
		return nil, errors.New("wrong naming convention: ip")
	}

	var nc *NamingConvention = &NamingConvention{
		Prefix:     parts[0],
		Ip:         ip.To4(),
		Mask:       (*ipNet).Mask,
		ObjectName: strings.Join(parts[1:len(parts)-5], "-"),
	}
	return nc, nil
}

func (nc *NamingConvention) Is(value string) bool {
	if nc == nil {
		return false
	}
	return nc.Prefix == value
}

func (nc *NamingConvention) IsVpcBridge() bool {
	return nc.Is("chbrvpc")
}

func (nc *NamingConvention) IsDefaultBridge() bool {
	return nc.Is("chbrdef")
}

func (nc *NamingConvention) IsVirtualMachineTap() bool {
	return nc.Is("chtap")
}

func NewVirtualMachineNetworkUtility(bridgeName string) *VirtualMachineNetworkUtility {
	return &VirtualMachineNetworkUtility{
		bridgeDefaultName: bridgeName,
	}
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
		nc, err := ParseDeviceName(link.Attrs().Name)
		if err != nil {
			return false
		}
		if !nc.IsVirtualMachineTap() {
			return false
		}
		return nc.ObjectName == vmId
	})
}

func (nm *VirtualMachineNetworkUtility) GetAllTapDevices() ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		tuntap, ok := link.(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			return false
		}
		nc, err := ParseDeviceName(link.Attrs().Name)
		if err != nil {
			return false
		}
		if !nc.IsVirtualMachineTap() {
			return false
		}
		if link.Attrs().MasterIndex == 0 {
			return false
		}
		master, err := netlink.LinkByIndex(link.Attrs().MasterIndex)
		if err != nil {
			return false
		}
		nc, err = ParseDeviceName(master.Attrs().Name)
		if err != nil {
			return false
		}
		if !(nc.IsDefaultBridge() || nc.IsVpcBridge()) {
			return false
		}
		return true
	})
}

func (nm *VirtualMachineNetworkUtility) GetAllVpcBridgeDevices() ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		_, ok := link.(*netlink.Bridge)
		if !ok {
			return false
		}
		nc, err := ParseDeviceName(link.Attrs().Name)
		if err != nil {
			return false
		}
		return nc.IsVpcBridge()
	})
}
