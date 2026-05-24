package networkvpc

import (
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteNetwork_Parse(t *testing.T) {
	row := []byte{DELETE_NETWORK}
	tenant := uuid.UUID{
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
	}
	row = append(row, tenant[:]...)
	row = append(row, []byte{
		byte(10),
		byte(0),
		byte(0),
		byte(0),
		byte(24),
	}...)
	deleteNetwork := &DeleteNetwork{}
	err := deleteNetwork.Parse(row, 0)
	assert.Nil(t, err, "No errors expected in Parse method")
	assert.Equal(t, tenant.String(), deleteNetwork.GetTenant(), "Expect to get the correct tenant id")
	assert.Equal(t, "10.0.0.0/24", deleteNetwork.GetNetworkString(), "Expected the correct network")
}

func Test_DeleteNetwork_Row(t *testing.T) {
	tenant := uuid.UUID{
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
	}
	deleteNetwork := &DeleteNetwork{
		action: DELETE_NETWORK,
		tenant: tenant,
		network: net.IPNet{
			IP:   net.IPv4(byte(10), byte(0), byte(0), byte(0)),
			Mask: net.IPv4Mask(byte(255), byte(255), byte(255), byte(0)),
		},
	}
	row := deleteNetwork.Row()
	expectedRow := []byte{DELETE_NETWORK}
	expectedRow = append(expectedRow, tenant[:]...)
	expectedRow = append(expectedRow, []byte{
		byte(10),
		byte(0),
		byte(0),
		byte(0),
		byte(24),
	}...)
	for i := 0; i < len(expectedRow); i++ {
		assert.Equal(t, expectedRow[i], row[i], "Expect the correct result for each byte in row")
	}
}
