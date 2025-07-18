package vmnetworking

import (
	"errors"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

type VirtualMachineNetworkingManager struct {
	bridgeDefaultName string
}

type NamingConvention struct {
	Prefix     string
	ObjectName string
	Ip         net.IP
}

func (nc *NamingConvention) Parse(name string) error {
	if nc == nil {
		return errors.New("uninitialized struct")
	}
	var parts []string = strings.Split(name, "-")
	if len(parts) < 6 {
		nc = nil
		return errors.New("wrong naming convention: format")
	}
	nc.Prefix = parts[0]
	var strIp string = strings.Join([]string{parts[len(parts)-4], parts[len(parts)-3], parts[len(parts)-2], parts[len(parts)-1]}, ".")
	var ip net.IP = net.ParseIP(strIp)
	if ip == nil {
		nc = nil
		return errors.New("wrong naming convention: ip")
	}
	nc.Ip = ip
	nc.ObjectName = strings.Join(parts[1:len(parts)-4], "-")
	return nil
}

func (nc *NamingConvention) Compare(value string) bool {
	if nc == nil {
		return false
	}
	if nc.Prefix == value {
		return true
	}
	return false
}

func (nc *NamingConvention) IsVpcBridge() bool {
	return nc.Compare("chbrvpc")
}

func (nc *NamingConvention) IsDefaultBridge() bool {
	return nc.Compare("chbrdef")
}

func (nc *NamingConvention) IsVirtualMachineTap() bool {
	return nc.Compare("chtap")
}

func NewVirtualMachineNetworkingManager(bridgeName string) *VirtualMachineNetworkingManager {
	return &VirtualMachineNetworkingManager{
		bridgeDefaultName: bridgeName,
	}
}

func (nm *VirtualMachineNetworkingManager) GetTapDevices(vmId string) ([]netlink.Link, error) {
	var netDevices []netlink.Link
	var err error
	var nc *NamingConvention = &NamingConvention{}
	var res []netlink.Link = []netlink.Link{}
	netDevices, err = nm.GetAllTapDevices()
	if err != nil {
		return []netlink.Link{}, err
	}
	for i := 0; i < len(netDevices); i++ {
		err = nc.Parse(netDevices[i].Attrs().Name)
		if err != nil {
			continue
		}
		if nc.ObjectName == vmId {
			res = append(res, netDevices[i])
		}
	}

	return res, nil
}

func (nm *VirtualMachineNetworkingManager) GetAllTapDevices() ([]netlink.Link, error) {
	var links []netlink.Link
	var err error
	var nc *NamingConvention = &NamingConvention{}
	var res []netlink.Link = []netlink.Link{}
	links, err = netlink.LinkList()
	if err != nil {
		return []netlink.Link{}, err
	}
	for i := 0; i < len(links); i++ {
		tuntap, ok := links[i].(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			continue
		}
		err = nc.Parse(netlink.NewLinkAttrs().Name)
		if err != nil {
			continue
		}
		if !nc.IsVirtualMachineTap() {
			continue
		}
		if links[i].Attrs().MasterIndex == 0 {
			continue
		}
		master, err := netlink.LinkByIndex(links[i].Attrs().MasterIndex)
		if err != nil {
			continue
		}
		err = nc.Parse(master.Attrs().Name)
		if err != nil {
			continue
		}
		if !(nc.IsDefaultBridge() || nc.IsVpcBridge()) {
			continue
		}
		res = append(res, links[i])
	}
	return res, nil
}

func (nm *VirtualMachineNetworkingManager) GetAllVpcBridgeDevices() ([]netlink.Link, error) {
	var links []netlink.Link
	var err error
	var nc *NamingConvention = &NamingConvention{}
	var res []netlink.Link = []netlink.Link{}
	links, err = netlink.LinkList()
	if err != nil {
		return []netlink.Link{}, err
	}
	for i := 0; i < len(links); i++ {
		_, ok := links[i].(*netlink.Bridge)
		if !ok {
			continue
		}
		err = nc.Parse(netlink.NewLinkAttrs().Name)
		if err != nil {
			continue
		}
		if !nc.IsVpcBridge() {
			continue
		}
		res = append(res, links[i])
	}
	return res, nil
}

type VirtualMachineNetworkingManagerService interface {
	GetTapDevices(vmId string) ([]netlink.Link, error)
	GetAllTapDevices() ([]netlink.Link, error)
	GetAllVpcBridgeDevices() ([]netlink.Link, error)
}
