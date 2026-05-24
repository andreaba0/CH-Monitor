package networkvpc

import (
	"errors"
	"net"
	"sync"
	"vmm/utils"

	"github.com/google/uuid"
)

const (
	ADD_NETWORK = iota
	DELETE_NETWORK
	DELETE_TENANT
)

type VpcWalStorage struct{}

func (s *VpcWalStorage) ReadSnapshot(path string) (map[string]map[string]string, error) {
	return utils.ReadGobFile[map[string]map[string]string](path)
}

func (s *VpcWalStorage) WriteSnapshot(path string, db map[string]map[string]string) error {
	return utils.WriteGobFile(path, db)
}

func (s *VpcWalStorage) CreateFile(path string) error {
	return utils.CreateFile(path)
}

func (s *VpcWalStorage) AppendRow(path string, row []byte) (int, error) {
	return utils.AppendOrCreateToFile(path, row)
}

func (s *VpcWalStorage) ReadChunk(path string, buffer []byte, index int64) (int, error) {
	return utils.ReadChunkFromFile(path, buffer, index)
}

type VpcWalStorageRepository interface {
	ReadSnapshot(path string) (map[string]map[string]string, error)
	WriteSnapshot(path string, db map[string]map[string]string) error
	CreateFile(path string) error
	AppendRow(path string, row []byte) (int, error)
	ReadChunk(path string, buffer []byte, index int64) (int, error)
}

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
	storage      VpcWalStorageRepository
}

func NewVpcManager(snapshotPath string, changesPath string) *VpcManager {
	return &VpcManager{
		snapshotPath: snapshotPath,
		changesPath:  changesPath,
		database:     make(map[string]map[string]string),
		storage:      new(VpcWalStorage),
	}
}

// When this method is called, AOF rows are loaded into main datastructure
// then the main data structure is stored back in filesystem and AOF file is cleared
func (vpcManager *VpcManager) LoadFromStorage(storagePath string) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	db, err := vpcManager.storage.ReadSnapshot(vpcManager.GetSnapshotFilePath())
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
	err := vpcManager.storage.WriteSnapshot(vpcManager.GetSnapshotFilePath(), vpcManager.database)
	if err != nil {
		return err
	}
	err = vpcManager.storage.CreateFile(vpcManager.GetLogFilePath())
	if err != nil {
		return err
	}
	return nil
}

func (VpcManager *VpcManager) processBuffer(index int, buffer []byte) (BlobData, error) {
	var data BlobData
	switch buffer[index] {
	case ADD_NETWORK:
		data = new(AddNetwork)
	case DELETE_NETWORK:
		data = new(DeleteNetwork)
	case DELETE_TENANT:
		data = new(DeleteTenant)
	default:
		return nil, errors.New("unknow data type")
	}
	err := data.Parse(buffer, index)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (vpcManager *VpcManager) readChangesFile() error {
	var index int64 = 0
	var err error
	var buffer []byte = make([]byte, 2048)
	for {
		n, err := vpcManager.storage.ReadChunk(vpcManager.changesPath, buffer, index)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		var localIndex int = 0
		for {
			data, err := vpcManager.processBuffer(localIndex, buffer)
			if err != nil {
				index += int64(localIndex)
				break
			}
			localIndex += data.GetRowSize()
			if obj, ok := data.(*AddNetwork); ok {
				vpcManager.addNetwork(obj, false)
				continue
			}
			if obj, ok := data.(*DeleteNetwork); ok {
				vpcManager.deleteNetwork(obj, false)
				continue
			}
			if obj, ok := data.(*DeleteTenant); ok {
				vpcManager.deleteTenant(obj, false)
				continue
			}
		}
		if !errors.Is(err, &ErrNotEnoughBytes{}) {
			break
		}
	}
	return err
}

func (vpcManager *VpcManager) addNetwork(an *AddNetwork, store bool) error {
	network := an.GetNetworkString()
	bridge := an.GetBridgeName()
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
	var err error = nil
	if store {
		_, err = vpcManager.storage.AppendRow(vpcManager.changesPath, an.Row())
	}
	return err
}

func (vpcManager *VpcManager) deleteNetwork(dn *DeleteNetwork, store bool) error {
	network := dn.GetNetworkString()
	tenant := dn.GetTenant()
	if _, ok := vpcManager.database[tenant]; ok {
		delete(vpcManager.database[tenant], network)
	}
	var err error = nil
	if store {
		_, err = vpcManager.storage.AppendRow(vpcManager.changesPath, dn.Row())
	}
	return err
}

func (vpcManager *VpcManager) deleteTenant(dt *DeleteTenant, store bool) error {
	tenant := dt.GetTenant()
	delete(vpcManager.database, tenant)
	var err error = nil
	if store {
		_, err = vpcManager.storage.AppendRow(vpcManager.changesPath, dt.Row())
	}
	return err
}

func (vpcManager *VpcManager) AddNetwork(tenant uuid.UUID, network net.IPNet, bridge string) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	addNetwork := NewAddNetwork(tenant, network, bridge)
	return vpcManager.addNetwork(addNetwork, true)
}

func (vpcManager *VpcManager) DeleteNetwork(tenant uuid.UUID, network net.IPNet) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	deleteNetwork := NewDeleteNetwork(tenant, network)
	return vpcManager.deleteNetwork(deleteNetwork, true)
}

func (vpcManager *VpcManager) DeleteTenant(tenant uuid.UUID) error {
	vpcManager.mu.Lock()
	defer vpcManager.mu.Unlock()
	deleteTenant := NewDeleteTenant(tenant)
	return vpcManager.deleteTenant(deleteTenant, true)
}
