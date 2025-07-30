package vmnetworking

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamingConvention_Parse(t *testing.T) {
	nc, err := ParseDeviceName("chtap-k8s-control-plane-192-168-0-4-24")
	assert.Nil(t, err, "Expect err to be null")
	assert.Equal(t, nc.ObjectName, "k8s-control-plane", "Expect the correct interface name")
	assert.Equal(t, nc.Ip.String(), "192.168.0.4", "Expect the correct ip")
	assert.Equal(t, nc.Mask.String(), "ffffff00", "Expect the correct ip mask")
	nc, err = ParseDeviceName("chtap---k8s-control-plane-192-168-0-4-24")
	assert.Nil(t, err, "Expect err to be null")
	assert.True(t, nc.IsVirtualMachineTap(), "Expect to be a vm tap interface")
	assert.Equal(t, nc.ObjectName, "--k8s-control-plane")
}

func Test_GenerateDeviceName(t *testing.T) {
	_, ipNet, _ := net.ParseCIDR("192.168.0.4/24")
	vmnu := &VirtualMachineNetworkUtility{
		defaultBridgeNetwork: *ipNet,
	}
	ni := &NetworkIdentifier{
		Ip:        net.IPv4(byte(192), byte(168), byte(0), byte(4)),
		Mask:      net.IPv4Mask(byte(255), byte(255), byte(255), byte(0)),
		Tenant:    "28eec4cc-75a3-4b34-aaf2-e4a8e6e71ed2",
		GuestName: "28eec4cc-75a3-ffff-ffff-e4a8e6e71ed2",
	}
	obj, err := vmnu.GenereateDeviceName(ni)
	assert.Nil(t, err, "No errors expected")
	assert.Equal(t, *(obj.Bridge), "chbrdef-192-168-0-0-24")
	assert.Equal(t, *(obj.Tap), "chtap-28eec4cc-75a3-ffff-ffff-e4a8e6e71ed2-192-168-0-4-24")

	ni.Ip = net.IPv4(byte(10), byte(0), byte(0), byte(4))
	obj, err = vmnu.GenereateDeviceName(ni)
	assert.Nil(t, err, "No errors expected")
	assert.Equal(t, *(obj.Bridge), "chbrvpc-28eec4cc-75a3-4b34-aaf2-e4a8e6e71ed2-10-0-0-0-24")
	assert.Equal(t, *(obj.Tap), "chtap-28eec4cc-75a3-4b34-aaf2-e4a8e6e71ed2-10-0-0-4-24")
}

func Test_ParseIp4ToNetDeviceString(t *testing.T) {
	str := ParseIp4ToNetDeviceString(net.IPv4(byte(192), byte(168), byte(0), byte(4)).To4(), net.IPv4Mask(byte(255), byte(255), byte(255), byte(0)))
	assert.Equal(t, str, "192-168-0-4-24")
}
