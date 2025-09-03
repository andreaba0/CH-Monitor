package networkvpc

import (
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	row := []byte{
		// Action
		ADD_NETWORK,
		// Tenant UUID
		0xc8,
		0xdd,
		0xee,
		0x26,
		0xbc,
		0x2b,
		0x46,
		0x5a,
		0x97,
		0x26,
		0x48,
		0x3f,
		0x11,
		0x7c,
		0xf6,
		0x18,
		// CIDR
		byte(10),
		byte(0),
		byte(0),
		byte(0),
		byte(24),
	}
	row = append(row, []byte("chbr-0000000001")...)
	addNetwork := &AddNetwork{}
	err := addNetwork.Parse(row, 0)
	assert.Nil(t, err, "No errors expected in Parse method")
	assert.Equal(t, "c8ddee26-bc2b-465a-9726-483f117cf618", addNetwork.GetTenant(), "Expect to get the correct tenant id")
	assert.Equal(t, "10.0.0.0/24", addNetwork.GetNetworkString(), "Expected the correct network")
	assert.Equal(t, "chbr-0000000001", addNetwork.GetBridgeName(), "Expected the correct bridge name")
}

func Test_Row(t *testing.T) {
	addNetwork := &AddNetwork{
		action: ADD_NETWORK,
		tenant: uuid.UUID{
			0xc8,
			0xdd,
			0xee,
			0x26,
			0xbc,
			0x2b,
			0x46,
			0x5a,
			0x97,
			0x26,
			0x48,
			0x3f,
			0x11,
			0x7c,
			0xf6,
			0x18,
		},
		bridge: "chbr-0000000001",
		network: net.IPNet{
			IP:   net.IPv4(byte(10), byte(0), byte(0), byte(0)),
			Mask: net.IPv4Mask(byte(255), byte(255), byte(255), byte(0)),
		},
	}
	row := addNetwork.Row()
	expectedRow := []byte{
		// Action
		ADD_NETWORK,
		// Tenant UUID
		0xc8,
		0xdd,
		0xee,
		0x26,
		0xbc,
		0x2b,
		0x46,
		0x5a,
		0x97,
		0x26,
		0x48,
		0x3f,
		0x11,
		0x7c,
		0xf6,
		0x18,
		// CIDR
		byte(10),
		byte(0),
		byte(0),
		byte(0),
		byte(24),
	}
	expectedRow = append(expectedRow, []byte("chbr-0000000001")...)
	for i := 0; i < len(expectedRow); i++ {
		assert.Equal(t, expectedRow[i], row[i], "Expect the correct result for each byte in row")
	}
}
