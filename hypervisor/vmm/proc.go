package vmm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	cloudhypervisor "vmm/cloud_hypervisor"
	"vmm/utils"
)

func parseProcFolder(cloudHypervisorPath string, pid int, procPath string, done chan<- *cloudhypervisor.CloudHypervisor) {
	// Read content of exe and cwd files to get process info
	// Step 1: check if the binary is related to Cloud-Hypervisor
	// Step 2: push process info (unix socket) to pool of vms to provision

	var cmdPath = filepath.Join(procPath, strconv.Itoa(pid), "/cmdline")
	var exePath = filepath.Join(procPath, strconv.Itoa(pid), "/exe")
	var err error
	fmt.Printf("Check %s, %s\n", exePath, cmdPath)

	contentCmd, err := os.ReadFile(cmdPath)
	if err != nil {
		done <- nil
		return
	}
	contentExe, err := os.Readlink(exePath)
	if err != nil {
		done <- nil
		return
	}

	if contentExe != cloudHypervisorPath {
		done <- nil
		return
	}

	var cmdLine = strings.Join(strings.Split(string(contentCmd), "\000"), " ")

	cmdLineParsing, err := utils.ParseCmdLine(cmdLine)
	if err != nil {
		done <- nil
		return
	}
	if cmdLineParsing.Binary != cloudHypervisorPath {
		done <- nil
		return
	}

	var cmdSocketParsing *utils.CmdSocketParsing = utils.NewCmdSocketParsing()
	err = cmdSocketParsing.Parse(cmdLineParsing.Args["--api-socket"].(string))
	if err != nil {
		done <- nil
		return
	}

	done <- cloudhypervisor.LoadRunningInstance(pid, cmdSocketParsing.Path)
}

func LoadProcessData(hypervisorBinaryPath string) ([]*cloudhypervisor.CloudHypervisor, error) {
	// List all processes in /proc

	var chanPool int = 20
	var procs chan *cloudhypervisor.CloudHypervisor = make(chan *cloudhypervisor.CloudHypervisor, chanPool)
	var index int
	var poolIndex int
	var procList []*cloudhypervisor.CloudHypervisor = make([]*cloudhypervisor.CloudHypervisor, 0)

	files, err := os.ReadDir("/proc")
	if err != nil {
		log.Fatalf("Unable to list folders in /proc; %s", err.Error())
	}
	index = 0
	poolIndex = 0
	var chInstance *cloudhypervisor.CloudHypervisor
	var i int
	for index < len(files) {

		// We are looking for process folder, not files
		if !files[index].IsDir() {
			index += 1
			continue
		}

		// This folder has no PID as its name, so we can skip it
		if !utils.IsUnicodeDigit(files[index].Name()) {
			index += 1
			continue
		}

		pid, err := strconv.Atoi(files[index].Name())
		if err != nil {
			index += 1
			continue
		}

		if strconv.Itoa(pid) != files[index].Name() {
			index += 1
			continue
		}

		if poolIndex < 20 {
			go parseProcFolder(hypervisorBinaryPath, pid, "/proc/", procs)
			index += 1
			poolIndex += 1
			continue
		}
		// Process chan return
		if poolIndex == 20 {
			chInstance = <-procs
			poolIndex -= 1
			if chInstance == nil {
				continue
			}
			procList = append(procList, chInstance)
		}
	}

	for i = 0; i < poolIndex; i++ {
		chInstance = <-procs
		if chInstance == nil {
			continue
		}
		procList = append(procList, chInstance)
	}

	return procList, nil
}
