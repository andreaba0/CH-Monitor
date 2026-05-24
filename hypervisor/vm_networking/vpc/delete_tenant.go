package networkvpc

import (
	"github.com/google/uuid"
)

type DeleteTenant struct {
	action byte
	tenant uuid.UUID
}

func NewDeleteTenant(tenant uuid.UUID) *DeleteTenant {
	return &DeleteTenant{
		action: DELETE_TENANT,
		tenant: tenant,
	}
}

func (deleteTenant *DeleteTenant) Parse(blob []byte, index int) error {
	if len(blob) < deleteTenant.GetRowSize() {
		return &ErrNotEnoughBytes{}
	}
	tenant, err := uuid.FromBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	deleteTenant.tenant = tenant
	deleteTenant.action = blob[0]
	return nil
}

func (deleteTenant *DeleteTenant) Row() []byte {
	res := []byte{}
	res = append(res, deleteTenant.action)
	res = append(res, deleteTenant.tenant[:]...)
	return res
}

func (deleteTenant *DeleteTenant) GetRowSize() int {
	return 1 + 16
}

func (deleteTenant *DeleteTenant) GetTenant() string {
	return deleteTenant.tenant.String()
}
