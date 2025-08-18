package metadata

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	vmstorage "vmm/storage"
)

type MetadataManager struct {
	tapCounter    int
	bridgeCounter int
	tapPrefix     string
	bridgePrefix  string
	snapshot_path string
	mu            sync.Mutex
}

func NewMetadataManager(snapshot_path string) (*MetadataManager, error) {
	manifest, err := vmstorage.ReadJson[*Manifest](snapshot_path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		return &MetadataManager{
			tapCounter:    manifest.TapCounter,
			bridgeCounter: manifest.BridgeCounter,
			tapPrefix:     manifest.TapPrefix,
			bridgePrefix:  manifest.BridgePrefix,
			snapshot_path: snapshot_path,
		}, nil
	}
	return &MetadataManager{
		tapCounter:    0,
		bridgeCounter: 0,
		tapPrefix:     "tapch",
		bridgePrefix:  "brvpc",
		snapshot_path: snapshot_path,
	}, nil
}

func (mm *MetadataManager) doSnapshot() error {
	manifest := Manifest{
		TapCounter:    mm.tapCounter,
		BridgeCounter: mm.bridgeCounter,
		TapPrefix:     mm.tapPrefix,
		BridgePrefix:  mm.bridgePrefix,
	}
	return vmstorage.WriteJson(mm.snapshot_path, &manifest)
}

func (mm *MetadataManager) MakeSnapshot() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.doSnapshot()
}

func (mm *MetadataManager) GetNewTapName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	tapName := fmt.Sprintf("%s-%s", mm.tapPrefix, strconv.Itoa(mm.tapCounter+1))
	err := mm.doSnapshot()
	if err != nil {
		return "", err
	}
	mm.tapCounter = mm.tapCounter + 1
	return tapName, nil
}

func (mm *MetadataManager) GetNewBridgeName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	bridgeName := fmt.Sprintf("%s-%s", mm.bridgePrefix, strconv.Itoa(mm.bridgeCounter+1))
	err := mm.doSnapshot()
	if err != nil {
		return "", err
	}
	mm.bridgeCounter = mm.bridgeCounter + 1
	return bridgeName, nil
}
