package networkvpc

import (
	"testing"
)

type MockedStorageVpc struct {
	Route         string
	testFramework *testing.T
}

func (ms *MockedStorageVpc) ReadSnapshot(path string) (map[string]map[string]string, error) {
	return nil, nil
}

func (ms *MockedStorageVpc) WriteSnapshot(path string, db map[string]map[string]string) error {
	return nil
}

func (ms *MockedStorageVpc) CreateFile(path string) error {
	return nil
}

func (ms *MockedStorageVpc) AppendRow(path string, row []byte) (int, error) {
	return 0, nil
}

func Test_AddNetwork(t *testing.T) {

}
