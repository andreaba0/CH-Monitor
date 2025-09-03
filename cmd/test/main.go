package main

import (
	"fmt"
	"os"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func main() {
	/*err := netlink.LinkAdd(&netlink.Tuntap{
		Mode: netlink.TUNTAP_MODE_TAP,
		LinkAttrs: netlink.LinkAttrs{
			Name: "tap0-vm",
		},
	})*/
	/*err := netlink.LinkAdd(&netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: "br0-vm",
		},
	})*/
	link, err := netlink.LinkByName("tap0-vm")
	if err != nil {
		fmt.Printf("Link not found, %T", err)
		os.Exit(1)
	}
	err = netlink.LinkDel(link)
	fmt.Printf("err: %T, %s\n", err, err.Error())
	if err == unix.EEXIST {
		fmt.Printf("Link already exists\n")
	}
	//isNonExistent := reflect.TypeOf(err) == reflect.TypeOf(netlink.LinkNotFoundError{})
	//fmt.Printf("bool: %t\n", isNonExistent)
	os.Exit(0)
}
