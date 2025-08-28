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
			tapPrefix:     "tpvm-",
			bridgePrefix:  "brvm-",
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

func (ne *NetworkEnumerator) TapName(number uint32) string {
	return fmt.Sprintf("%s%s", ne.tapPrefix, strconv.FormatUint(uint64(number), 10))
}

func (ne *NetworkEnumerator) BridgeName(number uint32) string {
	return fmt.Sprintf("%s%s", ne.bridgePrefix, strconv.FormatUint(uint64(number), 10))
}

func (mm *NetworkEnumerator) GetNewTapName() (string, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	tapName := mm.TapName(mm.tapCounter + 1)
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
	bridgeName := mm.BridgeName(mm.bridgeCounter + 1)
	err := mm.doSnapshot()
	if err != nil {
		return "", err
	}
	mm.bridgeCounter = mm.bridgeCounter + 1
	return bridgeName, nil
}
