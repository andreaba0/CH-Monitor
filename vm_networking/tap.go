package vmnetworking

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

func CreateTapInterface(tapName string, master netlink.Link) error {
	var err error
	var tap netlink.Tuntap = netlink.Tuntap{
		Mode:       netlink.TUNTAP_MODE_TAP,
		NonPersist: false,
		LinkAttrs: netlink.LinkAttrs{
			Name:        tapName,
			MasterIndex: master.Attrs().Index,
		},
	}
	err = netlink.LinkAdd(&tap)
	if err != nil {
		return err
	}
	return nil
}

func GenerateDefTapName(ip string, mask string) string {
	var ipStr string = strings.Join(strings.Split(ip, "."), "-")
	ipMask := net.IPMask(net.ParseIP(mask).To4())
	ones, _ := ipMask.Size()
	var name string = fmt.Sprintf("chtap-%s-%s", ipStr, strconv.Itoa(ones))
	return name
}
