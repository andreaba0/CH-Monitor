package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CmdLineParsing_Parse(t *testing.T) {
	var command = "/bin/cloud-hypervisor-static --api-socket path=/tmp/cloud-hypervisor/46c97539-797a-4cf6-b4b6-31f9909d9401.sock"
	var cmdLineParsing *CmdLineParsing = NewCmdLineParsing()
	err := cmdLineParsing.Parse(command)
	assert.Nil(t, err, "No errors should be raised")
	assert.Equal(t, cmdLineParsing.Binary, "/bin/cloud-hypervisor-static")
	assert.Equal(t, cmdLineParsing.Args["--api-socket"], "path=/tmp/cloud-hypervisor/46c97539-797a-4cf6-b4b6-31f9909d9401.sock")
}

func Test_CmdSocketParsing_Parse(t *testing.T) {
	var path = "path=/tmp/cloud-hypervisor/46c97539-797a-4cf6-b4b6-31f9909d9401.sock"
	var cmdSocketParsing *CmdSocketParsing = NewCmdSocketParsing()
	err := cmdSocketParsing.Parse(path)
	assert.Nil(t, err, "No errors should be raised")
	assert.Equal(t, cmdSocketParsing.Path, "/tmp/cloud-hypervisor/46c97539-797a-4cf6-b4b6-31f9909d9401.sock")
}
