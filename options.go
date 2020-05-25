package fb_nbackup

import (
	"io"
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
	WithTriggers(),
}

// WithCredential run a command with user privileges
func WithCredential(c *Credential) Option {
	return func(m *Manager) {
		m.credential = c
	}
}

// WithOutWriter run command with "out" output stream
func WithOutWriter(out io.Writer) Option {
	return func(m *Manager) {
		m.output = out
	}
}

// WithOutWriter run command with "out" error output stream
func WithErrWriter(out io.Writer) Option {
	return func(m *Manager) {
		m.outputErr = out
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

// WithoutTriggers Do not run databases triggers
func WithoutTriggers() Option {
	return func(m *Manager) {
		m.noDBTriggers = true
	}
}

// WithTriggers Run databases triggers
func WithTriggers() Option {
	return func(m *Manager) {
		m.noDBTriggers = false
	}
}

// WithCommandPath use alternative startup command nbackup
func WithCommandPath(s string) Option {
	return func(m *Manager) {
		m.command = s
	}
}
