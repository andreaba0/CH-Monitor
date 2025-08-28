package networkvpc

import (
	"errors"
	"net"
	"sync"
	storage "vmm/storage"

	"github.com/google/uuid"
)

const (
	EOF = iota
	ADD_NETWORK
	DELETE_NETWORK
	DELETE_TENANT
)

type BlobData interface {
	Parse(blob []byte, index int) error
	Row() []byte
	GetRowSize() int
}

type VpcManager struct {
	snapshotPath string
	changesPath  string
	database     map[string]map[string]string
	mu           sync.Mutex
}

func NewVpcManager(snapshotPath string, changesPath string) *VpcManager {
	return &VpcManager{
		snapshotPath: snapshotPath,
		changesPath:  changesPath,
		database:     make(map[string]map[string]string),
	}
}

// When this method is called, AOF rows are loaded into main datastructure
// then the main data structure is stored back in filesystem and AOF file is cleared
func (vpcManager *VpcManager) LoadFromStorage(storagePath string) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	db, err := storage.ReadJson[map[string]map[string]string](vpcManager.GetSnapshotFilePath())
	if err != nil {
		return err
	}
	vpcManager.database = db
	err = vpcManager.readChangesFile()
	if err != nil {
		return err
	}
	err = vpcManager.doSnapshot()
	if err != nil {
		return err
	}
	return nil
}

func (vpcManager *VpcManager) GetLogFilePath() string {
	return vpcManager.changesPath
}

func (vpcManager *VpcManager) GetSnapshotFilePath() string {
	return vpcManager.snapshotPath
}

// This method stores the main datastructure in filesystem and clears AOF file
func (vpcManager *VpcManager) doSnapshot() error {
	err := storage.WriteJson(vpcManager.GetSnapshotFilePath(), vpcManager.database)
	if err != nil {
		return err
	}
	err = storage.CreateFile(vpcManager.GetLogFilePath())
	if err != nil {
		return err
	}
	return nil
}

func (vpcManager *VpcManager) readChangesFile() error {
	var index int64 = 0
	cache := NewChunkCache(vpcManager.GetLogFilePath(), 512)
	for {
		buffer := cache.GetBuffered(index)
		if cache.BufferAndIndexAreAtEndOfFile(index) {
			break
		}
		data, err := vpcManager.getNextRow(buffer, index)
		if err != nil {
			if errors.Is(err, &ErrNotEnoughBytes{}) {
				cache.SlideBufferToIndex(index)
				continue
			} else {
				return err
			}
		}
		index += int64(data.GetRowSize())
		if obj, ok := data.(*AddNetwork); ok {
			vpcManager.addNetwork(obj, false)
			continue
		}
		return errors.New("unknow data type found")
	}
	return nil
}

func (vpcManager *VpcManager) getNextRow(buffer []byte, index int64) (BlobData, error) {
	localIndex := index % int64(len(buffer))
	var data BlobData
	switch buffer[localIndex] {
	case ADD_NETWORK:
		data = new(AddNetwork)
	case DELETE_NETWORK:
		data = new(DeleteNetwork)
	case DELETE_TENANT:
		data = new(DeleteTenant)
	default:
		return nil, errors.New("unknow data type")
	}
	err := data.Parse(buffer, int(localIndex))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (vpcManager *VpcManager) addNetwork(an *AddNetwork, store bool) error {
	network, err := an.GetNetworkString()
	if err != nil {
		return err
	}
	bridge := an.GetBridgeNumber()
	if _, ok := vpcManager.database[an.GetTenant()]; !ok {
		vpcManager.database[an.GetTenant()] = make(map[string]string)
	}
	if v, ok := vpcManager.database[an.GetTenant()][network]; ok {
		if v != bridge {
			return errors.New("mismatch in bridge name for a tenant")
		}
	} else {
		vpcManager.database[an.GetTenant()][network] = bridge
	}
	if store == true {
		_ = an.Row()
	}
	return nil
}

func (vpcManager *VpcManager) deleteNetwork(dn *DeleteNetwork, store bool) error {
	return nil
}

func (vpcManager *VpcManager) deleteTenant(dt *DeleteTenant, store bool) error {
	return nil
}

func (vpcManager *VpcManager) AddNetwork(tenant uuid.UUID, network net.IPNet, bridge string) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	addNetwork := NewAddNetwork(tenant, network, bridge)
	err := vpcManager.addNetwork(addNetwork, true)
	return err
}
