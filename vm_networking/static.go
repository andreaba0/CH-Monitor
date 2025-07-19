package vmnetworking

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

func CreateTapInterface(vmId string, networkAddress net.IP, master netlink.Link) error {
	var ipv4 net.IP = networkAddress.To4()
	var err error
	if ipv4 == nil {
		return errors.New("invalid ip address")
	}
	var tap netlink.Tuntap = netlink.Tuntap{
		Mode:       netlink.TUNTAP_MODE_TAP,
		NonPersist: false,
		LinkAttrs: netlink.LinkAttrs{
			Name: fmt.Sprintf(
				"chtap-%s-%s-%s-%s-%s",
				strings.ToLower(vmId),
				strconv.Itoa(int(ipv4[0])),
				strconv.Itoa(int(ipv4[1])),
				strconv.Itoa(int(ipv4[2])),
				strconv.Itoa(int(ipv4[3])),
			),
			MasterIndex: master.Attrs().Index,
		},
	}
	err = netlink.LinkAdd(&tap)
	if err != nil {
		return err
	}
	return nil
}

func CreateBridgeInterface(placeholder string, networkAddress net.IP, master *netlink.Link) error {
	var ipv4 net.IP = networkAddress.To4()
	var index int = 0
	var err error
	if ipv4 == nil {
		return errors.New("invalid ip address")
	}
	if master != nil {
		index = (*master).Attrs().Index
	}
	var bridge netlink.Bridge = netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: fmt.Sprintf(
				"chtap-%s-%s-%s-%s-%s",
				strings.ToLower(placeholder),
				strconv.Itoa(int(ipv4[0])),
				strconv.Itoa(int(ipv4[1])),
				strconv.Itoa(int(ipv4[2])),
				strconv.Itoa(int(ipv4[3])),
			),
			MasterIndex: index,
		},
	}
	err = netlink.LinkAdd(&bridge)
	if err != nil {
		return err
	}
	return nil
}
