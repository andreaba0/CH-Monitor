package cloudhypervisor

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
)

type CloudHypervisor struct {
	pid           int
	manifest      *Manifest
	binaryPath    string
	httpClient    *http.Client
	defaultBridge netlink.Link
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

func LoadFromExisting(pid int, socket string, manifest *Manifest, binaryPath string, defaultBridge netlink.Link) (*CloudHypervisor, error) {
	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		pid:           pid,
		manifest:      manifest,
		binaryPath:    binaryPath,
		httpClient:    CreateTransportSocket(socket),
		defaultBridge: defaultBridge,
	}
	return cloudHypervisor, nil
}

func NewCloudHypervisor(socket string, manifest *Manifest, binaryPath string, defaultBridge netlink.Link) (*CloudHypervisor, error) {

	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		manifest:      manifest,
		binaryPath:    binaryPath,
		httpClient:    CreateTransportSocket(socket),
		defaultBridge: defaultBridge,
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

type RunningInstance struct {
	socket          string
	cloudHypervisor *CloudHypervisor
}

func LoadRunning(binaryPath string) []RunningInstance {

}
