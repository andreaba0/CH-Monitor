package networkvpc

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteTenant_Parse(t *testing.T) {
	row := []byte{DELETE_TENANT}
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
	deleteTenant := &DeleteTenant{}
	err := deleteTenant.Parse(row, 0)
	assert.Nil(t, err, "No errors expected in Parse method")
	assert.Equal(t, tenant.String(), deleteTenant.GetTenant(), "Expect to get the correct tenant id")
}

func Test_DeleteTenant_Row(t *testing.T) {
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
	deleteTenant := &DeleteTenant{
		action: DELETE_NETWORK,
		tenant: tenant,
	}
	row := deleteTenant.Row()
	expectedRow := []byte{DELETE_NETWORK}
	expectedRow = append(expectedRow, tenant[:]...)
	for i := 0; i < len(expectedRow); i++ {
		assert.Equal(t, expectedRow[i], row[i], "Expect the correct result for each byte in row")
	}
}
