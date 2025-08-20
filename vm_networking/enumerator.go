package vmnetworking

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	vmstorage "vmm/storage"
)

type NetworkEnumerator struct {
	tapCounter    uint32
	bridgeCounter uint32
	tapPrefix     string
	bridgePrefix  string
	mu            sync.Mutex
	snapshot_path string
}

func NewNetworkEnumerator(snapshot_path string) (*NetworkEnumerator, error) {
	manifest, err := vmstorage.ReadJson[Manifest](snapshot_path)
	if err == nil {
		return &NetworkEnumerator{
			tapCounter:    manifest.TapCounter,
			bridgeCounter: manifest.BridgeCounter,
			tapPrefix:     manifest.TapPrefix,
			bridgePrefix:  manifest.BridgePrefix,
			snapshot_path: snapshot_path,
		}, nil
	}
	if os.IsNotExist(err) {
		return &NetworkEnumerator{
			tapCounter:    0,
			bridgeCounter: 0,
			tapPrefix:     "tapch",
			bridgePrefix:  "brvpc",
			snapshot_path: snapshot_path,
		}, nil
	}
	return nil, err
}

func (mm *NetworkEnumerator) doSnapshot() error {
	manifest := Manifest{
		TapCounter:    mm.tapCounter,
		BridgeCounter: mm.bridgeCounter,
		TapPrefix:     mm.tapPrefix,
		BridgePrefix:  mm.bridgePrefix,
	}
	return vmstorage.WriteJson(mm.snapshot_path, &manifest)
}

func (mm *NetworkEnumerator) MakeSnapshot() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	return mm.doSnapshot()
}

func (mm *NetworkEnumerator) GetNewTapName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	tapName := fmt.Sprintf("%s%s", mm.tapPrefix, strconv.FormatUint(uint64(mm.tapCounter+1), 10))
	err := mm.doSnapshot()
	if err != nil {
		return "", err
	}
	mm.tapCounter = mm.tapCounter + 1
	return tapName, nil
}

func (mm *NetworkEnumerator) GetNewBridgeName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	bridgeName := fmt.Sprintf("%s%s", mm.bridgePrefix, strconv.FormatUint(uint64(mm.bridgeCounter+1), 10))
	err := mm.doSnapshot()
	if err != nil {
		return "", err
	}
	mm.bridgeCounter = mm.bridgeCounter + 1
	return bridgeName, nil
}
