package vmnetworking

import (
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

type NetworkManager struct {
	defaultBridge netlink.Link
}

func NewNetworkManager(defaultBridge string) (*NetworkManager, error) {
	link, err := netlink.LinkByName(defaultBridge)
	if err != nil {
		return nil, err
	}
	return &NetworkManager{
		defaultBridge: link,
	}, nil
}

func (nm *NetworkManager) CreateTapInterface(ip net.IP, mask net.IPMask, tenant string) (netlink.Link, error) {
	ones, _ := mask.Size()
	ipStr := strings.Join(strings.Split(ip.String(), "."), "-")
	tapNameParts := []string{
		"tap",
		tenant,
		ipStr,
		strconv.Itoa(ones),
	}
	tapName := strings.Join(tapNameParts, "-")
	tap := netlink.Tuntap{
		Mode:       netlink.TUNTAP_MODE_TAP,
		NonPersist: false,
		LinkAttrs: netlink.LinkAttrs{
			Name: tapName,
		},
	}
	err := netlink.LinkAdd(&tap)
	if err != nil {
		return nil, err
	}
	return &tap, nil
}

func (nm *NetworkManager) CreateVpcBridge(network net.IPNet, tenant string) (netlink.Link, error) {
	return nil, nil
}

func (nm *NetworkManager) ConnectTapToVpc(tap netlink.Link, bridge netlink.Link) error {
	return netlink.LinkSetMaster(tap, bridge)
}

func (nm *NetworkManager) ConnectTapToDefault(tap netlink.Link) error {
	return netlink.LinkSetMaster(tap, nm.defaultBridge)
}
