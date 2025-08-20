package networkvpc

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/google/uuid"
)

type AddNetwork struct {
	action  byte
	tenant  uuid.UUID
	network net.IPNet
	bridge  uint32
	nextRow uint64
}

func NewAddNetwork(tenant uuid.UUID, network net.IPNet, bridge uint32) *AddNetwork {
	return &AddNetwork{
		action:  ADD_NETWORK,
		tenant:  tenant,
		network: network,
		bridge:  bridge,
	}
}

func (addNetwork *AddNetwork) Parse(blob []byte, index uint64) error {
	if uint64(len(blob)) < index+1+16+5+4 {
		return errors.New("not enough bytes")
	}
	tenant, err := uuid.ParseBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	addNetwork.tenant = tenant
	addNetwork.action = blob[0]
	addNetwork.nextRow = index + 1 + 16 + 5 + 4
	addNetwork.bridge = binary.BigEndian.Uint32(blob[index+1+16+5 : index+1+16+5+4])
	ip := net.IP(blob[index+1+16 : index+1+16+4])
	maskSize := int(blob[index+1+16+4])
	mask := net.CIDRMask(maskSize, 32)
	addNetwork.network = net.IPNet{
		IP:   ip.Mask(mask),
		Mask: mask,
	}
	return nil
}

func (addNetwork *AddNetwork) Row() []byte {
	res := []byte{}
	res = append(res, addNetwork.action)
	res = append(res, addNetwork.tenant[:]...)
	res = append(res, addNetwork.network.IP.To4()...)
	ones, _ := addNetwork.network.Mask.Size()
	res = append(res, byte(ones))
	res = append(res, byte(addNetwork.bridge))
	return res
}

func (addNetwork *AddNetwork) GetNextRow() uint64 {
	return addNetwork.nextRow
}
