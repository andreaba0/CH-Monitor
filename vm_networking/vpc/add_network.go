package networkvpc

import (
	"net"
	vmnetworking "vmm/vm_networking"

	"github.com/google/uuid"
)

type AddNetwork struct {
	action  byte
	tenant  uuid.UUID
	network net.IPNet
	bridge  string
}

func NewAddNetwork(tenant uuid.UUID, network net.IPNet, bridge string) *AddNetwork {
	return &AddNetwork{
		action:  ADD_NETWORK,
		tenant:  tenant,
		network: network,
		bridge:  bridge,
	}
}

func (addNetwork *AddNetwork) Parse(blob []byte, index int) error {
	if len(blob) < index+1+16+5+4 {
		return &ErrNotEnoughBytes{}
	}
	tenant, err := uuid.ParseBytes(blob[index+1 : index+1+16])
	if err != nil {
		return err
	}
	addNetwork.tenant = tenant
	addNetwork.action = blob[0]
	addNetwork.bridge = string(blob[index+1+16+5 : index+1+16+5+15])
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
	res = append(res, []byte(addNetwork.bridge)...)
	return res
}

func (addNetwork *AddNetwork) GetRowSize() int {
	return 1 + 16 + 5 + 15
}

func (addNetwork *AddNetwork) GetNetworkString() (string, error) {
	return vmnetworking.NetworkToCIDR(addNetwork.network)
}

func (addNetwork *AddNetwork) GetBridgeNumber() string {
	return addNetwork.bridge
}

func (addNetwork *AddNetwork) GetTenant() string {
	return addNetwork.tenant.String()
}
