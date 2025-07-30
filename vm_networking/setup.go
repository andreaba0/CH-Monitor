package vmnetworking

import (
	"errors"

	"github.com/vishvananda/netlink"
)

func SetupNetworking(defaultNetwork string) (*VirtualMachineNetworkUtility, error) {
	NICs, err := filterLinksBy(func(link netlink.Link) bool {
		if link.Attrs().PermHWAddr == nil || string(link.Attrs().PermHWAddr) == "" {
			return false
		}
		if link.Attrs().MasterIndex == 0 {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	if len(NICs) == 0 {
		return nil, errors.New("No NIC found")
	}
	var bridge netlink.Link = nil
	for i := 0; i < len(NICs); i++ {
		link, err := netlink.LinkByIndex(NICs[i].Attrs().MasterIndex)
		if err != nil {
			continue
		}
		if _, ok := link.(*netlink.Bridge); ok {
			bridge = link
			break
		}
	}
	if bridge == nil {
		return nil, errors.New("No bridge found")
	}

}
