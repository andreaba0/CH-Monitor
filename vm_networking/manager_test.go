package vmnetworking

import (
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
