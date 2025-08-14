package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/vishvananda/netlink"
)

func main() {
	_, err := netlink.LinkByName("unk")
	fmt.Printf("err: %s\n", err.Error())
	isNonExistent := reflect.TypeOf(err) == reflect.TypeOf(netlink.LinkNotFoundError{})
	fmt.Printf("bool: %t\n", isNonExistent)
	os.Exit(0)
}
