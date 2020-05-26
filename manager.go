package fb_nbackup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultCommand = "nbackup"
)

// Errors
var ErrUnknownArgumentType = errors.New("unknown argument type")

type Manager struct {
	command           string
	decompressCommand string
	direct            bool
	noDBTriggers      bool
	credential        *Credential
	executer          executer
	output            io.Writer
	outputErr         io.Writer
}

type executer interface {
	exec(ctx context.Context, commandLine string, args ...string) error
	execww(ctx context.Context, commandLine string, w io.Writer, args ...string) error
}

var (
	_ Backuper = (*Manager)(nil)
	_ executer = (*Manager)(nil)
)

func NewManager(opts ...Option) *Manager {
	manager := &Manager{}
	manager.executer = manager
	for _, option := range append(DefaultOptions, opts...) {
		option(manager)
	}
	return manager
}

func parseArguments(args ...interface{}) []string {
	parseArg := func(arg interface{}) []string {
		switch v := arg.(type) {
		case string:
			return []string{v}
		case []string:
			return v
		case Argument:
			return v.ToArgument()
		default:
			panic(fmt.Errorf("%w: %T", ErrUnknownArgumentType, arg))
		}
	}
	result := make([]string, 0)
	for _, arg := range args {
		result = append(result, parseArg(arg)...)
	}
	return result
}

func (m *Manager) buildCmd(args ...interface{}) (string, []string) {
	if m.credential != nil {
		args = append([]interface{}{m.credential}, args...)
	}
	argsParsed := parseArguments(args...)

	if len(m.command) == 0 {
		return defaultCommand, argsParsed
	}
	cmdParts := strings.Split(m.command, " ")
	return cmdParts[0], append(cmdParts[1:], argsParsed...)
}

func (m *Manager) exec(ctx context.Context, commandLine string,
	args ...string) error {
	cmd := exec.CommandContext(ctx, commandLine, args...)
	m.setWriter(cmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) execww(ctx context.Context, commandLine string, w io.Writer,
	args ...string) error {
	cmd := exec.CommandContext(ctx, commandLine, args...)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) setWriter(cmd *exec.Cmd) {
	if m.output != nil {
		cmd.Stdout = m.output
		cmd.Stderr = m.output
	}
}

// Regular expressions.
var (
	reVersion = regexp.MustCompile(`(?m)V(\d+\.)(\d+\.)(\d+\.)(\d+)`)
)

// Returns the nbachkup version.
// Will return an empty string if no version is found.
func (m *Manager) Version(ctx context.Context) (string, error) {
	cmd, args := m.buildCmd("-Z")
	data := &bytes.Buffer{}
	err := m.execww(ctx, cmd, data, args...)
	if err != nil {
		return "", err
	}
	for _, match := range reVersion.FindAllString(string(data.Bytes()), -1) {
		return match, nil
	}
	return "", nil
}

func (m *Manager) Lock(ctx context.Context, db string, returnSize bool) (int64, error) {
	commands := make([]interface{}, 0, 3)
	if returnSize {
		commands = append(commands, "-SIZE")
	}
	commands = append(commands, "-LOCK", db)

	cmd, args := m.buildCmd(commands...)
	data := &bytes.Buffer{}
	err := m.execww(ctx, cmd, data, args...)
	if err != nil {
		return 0, err
	}
	if !returnSize {
		return 0, nil
	}

	sData := strings.TrimSpace(string(data.Bytes()))

	size, err := strconv.Atoi(sData)
	if err != nil {
		return -1, err
	}
	return int64(size), nil
}

func (m *Manager) Unlock(ctx context.Context, db string) error {
	cmd, args := m.buildCmd("-UNLOCK", db)
	err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Fixup(ctx context.Context, db string) error {
	cmd, args := m.buildCmd("-FIXUP", db)
	err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Backup(ctx context.Context, level Level, db string, file string) error {
	cmd, args := m.buildCmd("-BACKUP", level, db, file)
	err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Restore(ctx context.Context, db string, files ...string) error {
	cmd, args := m.buildCmd("-RESTORE", db, files)
	err := m.exec(ctx, cmd, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) BackupTo(ctx context.Context, level Level, db string, w io.Writer) error {
	cmd, args := m.buildCmd("-BACKUP", level, db, "stdout")
	err := m.execww(ctx, cmd, w, args...)
	if err != nil {
		return err
	}
	return nil
}
