package vmnetworking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamingConvention_Parse(t *testing.T) {
	var nc *NamingConvention = &NamingConvention{}
	var err error = nc.Parse("chtap-k8s-control-plane-192-168-0-4")
	assert.Nil(t, err, "Expect err to be null")
	err = nc.Parse("chtap---k8s-control-plane-192-168-0-4")
	assert.Nil(t, err, "Expect err to be null")
	assert.True(t, nc.IsVirtualMachineTap(), "Expect to be a vm tap interface")
	assert.Equal(t, nc.ObjectName, "--k8s-control-plane")
}
