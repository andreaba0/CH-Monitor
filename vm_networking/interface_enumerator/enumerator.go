package interface_enumerator

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"vmm/utils"
)

type EnumeratorWalStorage struct{}

func (s *EnumeratorWalStorage) ReadSnapshot(path string) (*EnumeratorManifest, error) {
	return utils.ReadGobFile[*EnumeratorManifest](path)
}

func (s *EnumeratorWalStorage) WriteSnapshot(path string, db *EnumeratorManifest) error {
	return utils.WriteGobFile(path, db)
}

func (s *EnumeratorWalStorage) CreateFile(path string) error {
	return utils.CreateFile(path)
}

func (s *EnumeratorWalStorage) AppendRow(path string, row []byte) (int, error) {
	return utils.AppendOrCreateToFile(path, row)
}

func (s *EnumeratorWalStorage) ReadChunk(path string, buffer []byte, index int64) (int, error) {
	return utils.ReadChunkFromFile(path, buffer, index)
}

type EnumeratorWalStorageRepository interface {
	ReadSnapshot(path string) (*EnumeratorManifest, error)
	WriteSnapshot(path string, db *EnumeratorManifest) error
	CreateFile(path string) error
	AppendRow(path string, row []byte) (int, error)
	ReadChunk(path string, buffer []byte, index int64) (int, error)
}

type NetworkEnumerator struct {
	tapStorage    []bool
	tapIndex      int
	bridgeStorage []bool
	bridgeIndex   int
	tapPrefix     string
	bridgePrefix  string
	mu            sync.Mutex
	snapshot_path string
	storage       EnumeratorWalStorageRepository
}

func NewNetworkEnumerator(snapshot_path string) (*NetworkEnumerator, error) {
	storage := &EnumeratorWalStorage{}
	manifest, err := storage.ReadSnapshot(snapshot_path)
	if err == nil {
		return &NetworkEnumerator{
			tapStorage:    manifest.TapStorage,
			bridgeStorage: manifest.BridgeStorage,
			tapIndex:      0,
			bridgeIndex:   0,
			tapPrefix:     manifest.TapPrefix,
			bridgePrefix:  manifest.BridgePrefix,
			snapshot_path: snapshot_path,
			storage:       storage,
		}, nil
	}
	if os.IsNotExist(err) {
		return &NetworkEnumerator{
			tapStorage:    make([]bool, 512),
			tapIndex:      0,
			bridgeStorage: make([]bool, 512),
			bridgeIndex:   0,
			tapPrefix:     "tpvm-",
			bridgePrefix:  "brvm-",
			snapshot_path: snapshot_path,
			storage:       storage,
		}, nil
	}
	return nil, err
}

func (mm *NetworkEnumerator) doSnapshot() error {
	manifest := EnumeratorManifest{
		TapStorage:    mm.tapStorage,
		BridgeStorage: mm.bridgeStorage,
		TapPrefix:     mm.tapPrefix,
		BridgePrefix:  mm.bridgePrefix,
	}
	return mm.storage.WriteSnapshot(mm.snapshot_path, &manifest)
}

func (mm *NetworkEnumerator) MakeSnapshot() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.doSnapshot()
}

func (ne *NetworkEnumerator) TapName(number uint32) string {
	return fmt.Sprintf("%s%s", ne.tapPrefix, strconv.FormatUint(uint64(number), 10))
}

func (ne *NetworkEnumerator) BridgeName(number uint32) string {
	return fmt.Sprintf("%s%s", ne.bridgePrefix, strconv.FormatUint(uint64(number), 10))
}

func (mm *NetworkEnumerator) GenerateTapName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	index := -1
	for range len(mm.tapStorage) {
		if !mm.tapStorage[mm.tapIndex] {
			mm.tapIndex += 1
			if mm.tapIndex >= len(mm.tapStorage) {
				mm.tapIndex = 0
			}
		} else {
			index = mm.tapIndex
			break
		}
	}
	if index == -1 {
		return "", errors.New("no tap available")
	}
	tapName := mm.TapName(uint32(index))
	/*err := mm.doSnapshot()
	if err != nil {
		return "", err
	}*/
	return tapName, nil
}

func (mm *NetworkEnumerator) GenerateBridgeName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	index := -1
	for range len(mm.bridgeStorage) {
		if !mm.bridgeStorage[mm.bridgeIndex] {
			mm.bridgeIndex += 1
			if mm.bridgeIndex >= len(mm.bridgeStorage) {
				mm.bridgeIndex = 0
			}
		} else {
			index = mm.bridgeIndex
			break
		}
	}
	bridgeName := mm.BridgeName(uint32(index))
	/*err := mm.doSnapshot()
	if err != nil {
		return "", err
	}*/
	return bridgeName, nil
}
