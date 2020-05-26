package fb_nbackup

import (
	"io"
	"os"
)

// Options is the type for NBACKUP options.
type Option func(m *Manager)

// Credential used as an option to run nbackup from the user
type Credential struct {
	User             string
	Role             string
	Password         string
	PasswordFromFile string
}

var _ Argument = (*Credential)(nil)

// ToArgument implements Argument
func (c *Credential) ToArgument() []string {
	var args []string
	addArgument := func(name, arg string) {
		if arg != "" {
			args = append(args, name, arg)
		}
	}
	addArgument("-USER", c.User)
	addArgument("-ROLE", c.Role)
	addArgument("-PASSWORD", c.Password)
	addArgument("-FETCH_PASSWORD", c.PasswordFromFile)
	return args
}

// Nbackup default options
var DefaultOptions = []Option{
	WithDirect(false),
	WithoutTriggers(true),
	WithWriter(os.Stdout),
}

// WithCredential run a command with user privileges
func WithCredential(c *Credential) Option {
	return func(m *Manager) {
		m.credential = c
	}
}

// WithWriter  run command with "out" output stream
func WithWriter(out io.Writer) Option {
	return func(m *Manager) {
		m.output = out
	}
}

// WithDecompressCommand - Command to extract archives during restore
func WithDecompressCommand(command string) Option {
	return func(m *Manager) {
		m.decompressCommand = command
	}
}

// WithDirect - Use or not direct I/O when backing up database
func WithDirect(use bool) Option {
	return func(m *Manager) {
		m.direct = use
	}
}

// WithoutTriggers Run databases triggers
func WithoutTriggers(trigger bool) Option {
	return func(m *Manager) {
		m.noDBTriggers = trigger
	}
}

// WithCommandPath use alternative startup command nbackup
func WithCommandPath(s string) Option {
	return func(m *Manager) {
		m.command = s
	}
}
