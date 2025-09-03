package vmnetworking

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/vishvananda/netlink"
)

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

func GetTapDevices(vmId string) ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		tuntap, ok := link.(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			return false
		}
		return true
	})
}

func GetAllTapDevices() ([]netlink.Link, error) {
	return filterLinksBy(func(link netlink.Link) bool {
		tuntap, ok := link.(*netlink.Tuntap)
		if !ok || tuntap.Mode != netlink.TUNTAP_MODE_TAP {
			return false
		}
		if link.Attrs().MasterIndex == 0 {
			return false
		}
		_, err := netlink.LinkByIndex(link.Attrs().MasterIndex)
		return err == nil
	})
}

func ParseCIDR4(ipStr string, maskStr string) (net.IP, *net.IPNet, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil || ip.To4() == nil {
		return nil, nil, errors.New("expected an ipv4 address")
	}
	ipMask := net.ParseIP(maskStr)
	if ipMask == nil || ipMask.To4() == nil {
		return nil, nil, errors.New("expected an ipv4 mask")
	}
	mask := net.IPMask(ipMask.To4())
	if mask == nil {
		return nil, nil, errors.New("expected an ipv4 mask")
	}
	ones, _ := mask.Size()
	return net.ParseCIDR(fmt.Sprintf("%s/%s", ip.To4().String(), strconv.Itoa(ones)))
}

func NetworkToCIDR(network net.IPNet) (string, error) {
	return "", nil
}

func CreateTapDevice(name string, master netlink.Link) error {
	err := netlink.LinkAdd(&netlink.Tuntap{
		Mode: netlink.TUNTAP_MODE_TAP,
		LinkAttrs: netlink.LinkAttrs{
			Name:        name,
			MasterIndex: master.Attrs().Index,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateBridgeDevice(name string) error {
	err := netlink.LinkAdd(&netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
