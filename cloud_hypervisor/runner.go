package cloudhypervisor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type CloudHypervisor struct {
	pid        int
	HttpClient *http.Client
	RestServer *HypervisorRestServer
}

func waitSocketFileCreation(socket string) error {
	for i := 0; i < 50; i++ {
		_, err := os.Stat(socket)
		if err != nil && os.IsNotExist(err) {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("wait for socket creation: time expired")
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

func NewCloudHypervisor(binaryPath string, remoteUri string) (*CloudHypervisor, error) {
	var err error
	socketUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.New("error while generating uuid for socket file")
	}
	var socketPath string = fmt.Sprintf("/tmp/vm-net-%s.sock", socketUuid)

	cmd := exec.Command(binaryPath, "--api-socket", fmt.Sprintf("path=%s", socketPath))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	err = waitSocketFileCreation(socketPath)
	if err != nil {
		return nil, err
	}

	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		pid:        cmd.Process.Pid,
		HttpClient: CreateTransportSocket(socketPath),
		RestServer: NewHypervisorRestServer(remoteUri),
	}

	return cloudHypervisor, nil
}

func (ch *CloudHypervisor) Kill() error {

	// TODO:
	// 1. Terminate current cloud-hypervisor instance
	// 2. Drop all tap interfaces connected to the instance

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

func LoadRunningInstance(pid int, socketPath string) *CloudHypervisor {
	var cloudHypervisor *CloudHypervisor = &CloudHypervisor{
		pid:        pid,
		HttpClient: CreateTransportSocket(socketPath),
	}
	return cloudHypervisor
}

func (ch *CloudHypervisor) IsRunning() bool {
	return ch.pid > 0
}
