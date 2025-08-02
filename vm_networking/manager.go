package vmnetworking

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

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

func ParseIp4ToNetDeviceString(ip net.IP, mask net.IPMask) string {
	ones, _ := mask.Size()
	var parts []string = []string{}
	parts = append(parts, strconv.Itoa(int(ip[0])))
	parts = append(parts, strconv.Itoa(int(ip[1])))
	parts = append(parts, strconv.Itoa(int(ip[2])))
	parts = append(parts, strconv.Itoa(int(ip[3])))
	parts = append(parts, strconv.Itoa(ones))
	return strings.Join(parts, "-")
}

func (vmn *VirtualMachineNetworkUtility) GenereateDeviceName(ni *NetworkIdentifier) (*NetworkDeviceNameGeneration, error) {
	ones1, bits1 := vmn.defaultBridgeNetwork.Mask.Size()
	ones2, bits2 := ni.Mask.Size()
	if ones1 == ones2 && bits1 == bits2 && vmn.defaultBridgeNetwork.Contains(ni.Ip) {
		// IP is inside the default host network
		var bridge []string = []string{}
		bridge = append(bridge, "chbrdef")
		bridge = append(bridge, ParseIp4ToNetDeviceString(vmn.defaultBridgeNetwork.IP.To4(), ni.Mask))
		var bridgeStr string = strings.Join(bridge, "-")
		var tap []string = []string{}
		tap = append(tap, "chtap")
		tap = append(tap, ni.GuestName)
		tap = append(tap, ParseIp4ToNetDeviceString(ni.Ip.To4(), ni.Mask))
		var tapStr string = strings.Join(tap, "-")
		return &NetworkDeviceNameGeneration{
			Bridge: &bridgeStr,
			Tap:    &tapStr,
			VxLan:  nil,
		}, nil
	}
	var bridge []string = []string{}
	var networkIp = ni.Ip.To4().Mask(ni.Mask)
	bridge = append(bridge, "chbrvpc")
	bridge = append(bridge, ni.Tenant)
	bridge = append(bridge, ParseIp4ToNetDeviceString(networkIp.To4(), ni.Mask))
	var bridgeStr string = strings.Join(bridge, "-")
	var tap []string = []string{}
	tap = append(tap, "chtap")
	tap = append(tap, ni.Tenant)
	tap = append(tap, ParseIp4ToNetDeviceString(ni.Ip.To4(), ni.Mask))
	var tapStr string = strings.Join(tap, "-")
	return &NetworkDeviceNameGeneration{
		Bridge: &bridgeStr,
		Tap:    &tapStr,
		VxLan:  nil,
	}, nil
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
		_, err = ParseDeviceName(master.Attrs().Name)
		return err == nil
	})
}

func GenerateTapName(ip net.IP, ipNet net.IPNet) string {
	var str string = ParseIp4ToNetDeviceString(ip.To4(), ipNet.Mask)
	return fmt.Sprintf("chtap-%s", str)
}
