package vmnetworking

import (
	"net"
	"reflect"
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

func (nm *NetworkManager) GenerateTapName(ip net.IP, mask net.IPMask, tenant string) string {
	ones, _ := mask.Size()
	ipStr := strings.Join(strings.Split(ip.String(), "."), "-")
	tapNameParts := []string{
		"tap",
		tenant,
		ipStr,
		strconv.Itoa(ones),
	}
	tapName := strings.Join(tapNameParts, "-")
	return tapName
}

func (nm *NetworkManager) GetTapInterface(ip net.IP, mask net.IPMask, tenant string) (netlink.Link, error) {
	tapName := nm.GenerateTapName(ip, mask, tenant)
	return netlink.LinkByName(tapName)
}

func (nm *NetworkManager) GenerateVpcName(network net.IPNet, tenant string) string {
	ones, _ := network.Mask.Size()
	bridgeNameParts := []string{
		"br",
		tenant,
		strings.Join(strings.Split(network.IP.To4().String(), "."), "-"),
		strconv.Itoa(ones),
	}
	return strings.Join(bridgeNameParts, "-")
}

func (nm *NetworkManager) ConnectTapToVpc(tap netlink.Link, bridge netlink.Link) error {
	return netlink.LinkSetMaster(tap, bridge)
}

func (nm *NetworkManager) ConnectTapToDefault(tap netlink.Link) error {
	return netlink.LinkSetMaster(tap, nm.defaultBridge)
}

func (nm *NetworkManager) GetAndCreateIfNotExistsVpc(network net.IPNet, tenant string) (netlink.Link, error) {
	link, err := netlink.LinkByName(nm.GenerateVpcName(network, tenant))
	if err == nil {
		return link, nil
	}
	if reflect.TypeOf(err) != reflect.TypeOf(netlink.LinkNotFoundError{}) {
		return nil, err
	}
	link = &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: nm.GenerateVpcName(network, tenant),
		},
	}
	err = netlink.LinkAdd(link)
	if err != nil {
		return nil, err
	}
	return link, nil
}
