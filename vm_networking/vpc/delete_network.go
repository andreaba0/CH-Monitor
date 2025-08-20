package networkvpc

import (
	"errors"
	"net"

	"github.com/google/uuid"
)

type DeleteNetwork struct {
	action  byte
	tenant  uuid.UUID
	network net.IPNet
	nextRow uint64
}

func NewDeleteNetwork(tenant uuid.UUID, network net.IPNet) *DeleteNetwork {
	return &DeleteNetwork{
		action:  DELETE_NETWORK,
		tenant:  tenant,
		network: network,
	}
}

func (deleteNetwork *DeleteNetwork) Parse(blob []byte, index uint64) error {
	if uint64(len(blob)) < index+1+16+5 {
		return errors.New("not enough bytes")
	}
	tenant, err := uuid.ParseBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	deleteNetwork.tenant = tenant
	deleteNetwork.action = blob[0]
	deleteNetwork.nextRow = index + 1 + 16 + 5
	ip := net.IP(blob[index+1+16 : index+1+16+4])
	maskSize := int(blob[index+1+16+4])
	mask := net.CIDRMask(maskSize, 32)
	deleteNetwork.network = net.IPNet{
		IP:   ip.Mask(mask),
		Mask: mask,
	}
	return nil
}

func (deleteNetwork *DeleteNetwork) Row() []byte {
	res := []byte{}
	res = append(res, deleteNetwork.action)
	res = append(res, deleteNetwork.tenant[:]...)
	res = append(res, deleteNetwork.network.IP.To4()...)
	ones, _ := deleteNetwork.network.Mask.Size()
	res = append(res, byte(ones))
	return res
}

func (deleteNetwork *DeleteNetwork) GetNextRow() uint64 {
	return deleteNetwork.nextRow
}
