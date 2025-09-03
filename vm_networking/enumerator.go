package vmnetworking

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type storageEnumerator struct{}

func (s *storageEnumerator) ReadSnapshot(path string) (*EnumeratorManifest, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	var res *EnumeratorManifest
	decoder := gob.NewDecoder(fd)
	err = decoder.Decode(res)
	return res, err
}

func (s *storageEnumerator) WriteSnapshot(path string, manifest *EnumeratorManifest) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer fd.Close()
	encoder := gob.NewEncoder(fd)
	return encoder.Encode(manifest)
}

type storageEnumeratorService interface {
	ReadSnapshot(path string) (*EnumeratorManifest, error)
	WriteSnapshot(path string, manifest *EnumeratorManifest) error
}

type NetworkEnumerator struct {
	tapCounter    uint32
	bridgeCounter uint32
	tapPrefix     string
	bridgePrefix  string
	mu            sync.Mutex
	snapshot_path string
	storage       storageEnumeratorService
}

func NewNetworkEnumerator(snapshot_path string) (*NetworkEnumerator, error) {
	storage := new(storageEnumerator)
	manifest, err := storage.ReadSnapshot(snapshot_path)
	if err == nil {
		return &NetworkEnumerator{
			tapCounter:    manifest.TapCounter,
			bridgeCounter: manifest.BridgeCounter,
			tapPrefix:     manifest.TapPrefix,
			bridgePrefix:  manifest.BridgePrefix,
			snapshot_path: snapshot_path,
			storage:       storage,
		}, nil
	}
	if os.IsNotExist(err) {
		return &NetworkEnumerator{
			tapCounter:    0,
			bridgeCounter: 0,
			tapPrefix:     "tpvm-",
			bridgePrefix:  "brvm-",
			snapshot_path: snapshot_path,
			storage:       new(storageEnumerator),
		}, nil
	}
	return nil, err
}

func (mm *NetworkEnumerator) doSnapshot() error {
	manifest := EnumeratorManifest{
		TapCounter:    mm.tapCounter,
		BridgeCounter: mm.bridgeCounter,
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
