package vmnetworking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseCIDR4(t *testing.T) {
	ip, ipNet, err := ParseCIDR4("10.0.0.2", "255.255.255.0")
	assert.Nil(t, err, "Expect err to be null")
	assert.Equal(t, ip.To4().String(), "10.0.0.2", "Expect the correct ip address")
	ones, _ := ipNet.Mask.Size()
	assert.Equal(t, ones, 24, "Expect the correct subnet mask")
}
