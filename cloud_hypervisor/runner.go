package cloudhypervisor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/vishvananda/netlink"
)

type CloudHypervisor struct {
	pid           int
	manifest      *Manifest
	binaryPath    string
	httpClient    *http.Client
	defaultBridge netlink.Link
	socketPath    string
}

func CreateTransportSocket(socket string) *http.Client {
	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return net.Dial("unix", socket)
	}
	transport := &http.Transport{
		DialContext: dialer,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
	return client
}

func NewCloudHypervisor(manifest *Manifest, binaryPath string, defaultBridge netlink.Link) (*CloudHypervisor, error) {
	socketUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.New("error while generating uuid for socket file")
	}
	var socketPath string = fmt.Sprintf("/tmp/vm-net-%s.sock", socketUuid)
	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		manifest:      manifest,
		binaryPath:    binaryPath,
		httpClient:    CreateTransportSocket(socketPath),
		defaultBridge: defaultBridge,
		socketPath:    socketPath,
	}
	return cloudHypervisor, nil
}

func (ch *CloudHypervisor) Kill() error {
	proc, err := os.FindProcess(ch.pid)
	if err != nil {
		return errors.New("there was an error searching process by pid")
	}
	err = proc.Signal(syscall.SIGKILL)
	if err != nil {
		return errors.New("there was an error killing running process")
	}
	return nil
}

func LoadRunning(binaryPath string) []CloudHypervisor {
	return []CloudHypervisor{}
}
