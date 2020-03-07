package fb_nbackup

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	defaultCommand = "nbackup"
)

type Manager struct {
	command           string
	decompressCommand string
	direct            bool
	noDBTriggers      bool
	credintial        Credintial

	executer executer
}

type executer func(cmd string, args ...string) ([]byte, error)

func NewManager(opts ...Option) *Manager {
	manager := &Manager{}
	manager.executer = manager.exec
	for _, option := range opts {
		option(manager)
	}
	return manager
}

func (m *Manager) buildCmd(args ...string) (string, []string) {
	if len(m.command) == 0 {
		return defaultCommand, args
	}
	cmdParts := strings.Split(m.command, " ")
	return cmdParts[0], append(cmdParts[1:], args...)
}

func (m *Manager) exec(commandLine string, args ...string) ([]byte, error) {
	cmd := exec.Command(commandLine, args...)
	var bufOut, bufErr bytes.Buffer
	cmd.Stderr = &bufErr
	cmd.Stdout = &bufOut
	err := cmd.Run()
	if err != nil {
		fmt.Println(bufErr.String())
		return nil, err
	}
	return bytes.Join([][]byte{bufOut.Bytes(), bufErr.Bytes()}, nil), nil
}

var re = regexp.MustCompile(`(?m)V(\d+\.)(\d+\.)(\d+\.)(\d+)`)

func (m *Manager) Version() (string, error) {
	cmd, args := m.buildCmd("-Z")
	data, err := m.exec(cmd, args...)
	if err != nil {
		return "", err
	}
	for _, match := range re.FindAllString(string(data), -1) {
		return match, nil
	}
	return string(data), nil
}
