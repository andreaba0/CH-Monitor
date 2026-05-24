package networkvpc

import (
	"net"
	vmnetworking "vmm/vm_networking"

	"github.com/google/uuid"
)

type DeleteNetwork struct {
	action  byte
	tenant  uuid.UUID
	network net.IPNet
}

func NewDeleteNetwork(tenant uuid.UUID, network net.IPNet) *DeleteNetwork {
	return &DeleteNetwork{
		action:  DELETE_NETWORK,
		tenant:  tenant,
		network: network,
	}
}

func (deleteNetwork *DeleteNetwork) Parse(blob []byte, index int) error {
	if len(blob) < deleteNetwork.GetRowSize() {
		return &ErrNotEnoughBytes{}
	}
	tenant, err := uuid.FromBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	deleteNetwork.tenant = tenant
	deleteNetwork.action = blob[0]
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

func (deleteNetwork *DeleteNetwork) GetRowSize() int {
	return 1 + 16 + 5
}

func (deleteNetwork *DeleteNetwork) GetNetworkString() string {
	return vmnetworking.NetworkToCIDR4(deleteNetwork.network)
}

func (deleteNetwork *DeleteNetwork) GetTenant() string {
	return deleteNetwork.tenant.String()
}
