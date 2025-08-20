package networkvpc

import (
	"errors"
	"path/filepath"
	"sync"
	storage "vmm/storage"
)

const (
	ADD_NETWORK = iota
	DELETE_NETWORK
	DELETE_TENANT
)

type BlobData interface {
	Parse(blob []byte, index uint64) error
	Row() []byte
	GetNextRow() uint64
}

type VpcManager struct {
	storagePath string
	database    map[string]map[string]string
	mu          sync.Mutex
}

func NewVpcManager(storagePath string) *VpcManager {
	return &VpcManager{
		storagePath: storagePath,
		database:    make(map[string]map[string]string),
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
	return filepath.Join(vpcManager.storagePath, "changes.aof")
}

func (vpcManager *VpcManager) GetSnapshotFilePath() string {
	return filepath.Join(vpcManager.storagePath, "registry.bin")
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
	var blobSize int = 256
	for {
		blob, err := storage.ReadFileChunk(vpcManager.GetLogFilePath(), index, blobSize)
		if err != nil {
			return err
		}
		var data BlobData
		switch blob[0] {
		case ADD_NETWORK:
			data = new(AddNetwork)
		case DELETE_NETWORK:
			data = new(DeleteNetwork)
		case DELETE_TENANT:
			data = new(DeleteTenant)
		default:
			data = nil
		}
		if data == nil {
			return errors.New("unknow data type found")
		}
		err = data.Parse(blob, uint64(index))
		if err != nil {
			return err
		}

	}
	return nil
}
