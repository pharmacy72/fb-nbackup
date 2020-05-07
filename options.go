package fb_nbackup

import "os"

type Option func(m *Manager)

type Credential struct {
	User             string
	Role             string
	Password         string
	PasswordFromFile string
}

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

var DefaultOptions = []Option{
	WithDirect(false),
	WithTriggers(),
}

func WithCredential(c *Credential) Option {
	return func(m *Manager) {
		m.credential = c
	}
}

func WithStdOutput() Option {
	return func(m *Manager) {
		m.output = os.Stdout
		m.outputErr = os.Stderr
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

func WithoutTriggers() Option {
	return func(m *Manager) {
		m.noDBTriggers = true
	}
}

func WithTriggers() Option {
	return func(m *Manager) {
		m.noDBTriggers = false
	}
}

func WithCommandPath(s string) Option {
	return func(m *Manager) {
		m.command = s
	}
}
