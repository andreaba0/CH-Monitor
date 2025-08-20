package networkvpc

import (
	"errors"

	"github.com/google/uuid"
)

type DeleteTenant struct {
	action  byte
	tenant  uuid.UUID
	nextRow uint64
}

func NewDeleteTenant(tenant uuid.UUID) *DeleteNetwork {
	return &DeleteNetwork{
		action: DELETE_TENANT,
		tenant: tenant,
	}
}

func (deleteTenant *DeleteTenant) Parse(blob []byte, index uint64) error {
	if uint64(len(blob)) < index+1+16+5+4 {
		return errors.New("not enough bytes")
	}
	tenant, err := uuid.ParseBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	deleteTenant.tenant = tenant
	deleteTenant.action = blob[0]
	deleteTenant.nextRow = index + 1 + 16
	return nil
}

func (deleteTenant *DeleteTenant) Row() []byte {
	res := []byte{}
	res = append(res, deleteTenant.action)
	res = append(res, deleteTenant.tenant[:]...)
	return res
}

func (deleteTenant *DeleteTenant) GetNextRow() uint64 {
	return deleteTenant.nextRow
}
