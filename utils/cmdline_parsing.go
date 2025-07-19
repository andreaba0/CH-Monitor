package utils

import (
	"errors"
	"strings"

	"github.com/google/shlex"
)

type CmdLineParsing struct {
	Args   map[string]interface{}
	Binary string
}

func NewCmdLineParsing() *CmdLineParsing {
	return &CmdLineParsing{
		Args:   make(map[string]interface{}),
		Binary: "",
	}
}

func ParseCmdLine(command string) (*CmdLineParsing, error) {
	var cmdLine *CmdLineParsing = &CmdLineParsing{
		Binary: "",
		Args:   make(map[string]interface{}),
	}
	args, err := shlex.Split(command)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(args); i++ {
		if i == 0 {
			cmdLine.Binary = args[0]
			continue
		}
		if args[i][:2] == "--" && i == len(args)-1 {
			cmdLine.Args[args[i]] = true
			continue
		}
		if args[i][:2] != "--" {
			return nil, errors.New("expected a flag")
		}
		if args[i][:2] == "--" && args[i+1][:2] == "--" {
			cmdLine.Args[args[i]] = true
			continue
		}
		cmdLine.Args[args[i]] = args[i+1]
		i += 1
	}
	return cmdLine, nil
}

type CmdSocketParsing struct {
	Path string
}

func NewCmdSocketParsing() *CmdSocketParsing {
	return &CmdSocketParsing{
		Path: "",
	}
}

func (current *CmdSocketParsing) Parse(pathArg string) error {
	var parts []string = strings.Split(pathArg, "=")
	if len(parts) != 2 || parts[0] != "path" {
		return errors.New("unknown property")
	}

	current.Path = parts[1]

	return nil
}
