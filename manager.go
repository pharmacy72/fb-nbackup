package fb_nbackup

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
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

type executer func(context.Context, string, ...string) ([]byte, error)

var _ Backuper = (*Manager)(nil)

//TODO: Stderr -> error

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

func (m *Manager) exec(ctx context.Context, commandLine string,
	args ...string) ([]byte, error) {
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

func (m *Manager) Version(ctx context.Context) (string, error) {
	cmd, args := m.buildCmd("-Z")
	data, err := m.exec(ctx, cmd, args...)
	if err != nil {
		return "", err
	}
	for _, match := range re.FindAllString(string(data), -1) {
		return match, nil
	}
	return string(data), nil
}

func (m *Manager) Lock(ctx context.Context, db string, returnSize bool) (int64, error) {
	commands := make([]string, 0, 3)
	if returnSize {
		commands = append(commands, "-SIZE")
	}
	commands = append(commands, "-LOCK", db)

	cmd, args := m.buildCmd(commands...)
	data, err := m.exec(ctx, cmd, args...)
	if err != nil {
		return 0, err
	}
	if !returnSize {
		return 0, nil
	}

	sData := strings.TrimSpace(string(data))

	size, err := strconv.Atoi(sData)
	if err != nil {
		return -1, err
	}
	return int64(size), nil
}

func (m *Manager) Unlock(ctx context.Context, db string) error {
	cmd, args := m.buildCmd("-UNLOCK", db)
	_, err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Fixup(ctx context.Context, db string) error {
	cmd, args := m.buildCmd("-FIXUP", db)
	_, err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}
