package fb_nbackup

import "os"

type Option func(m *Manager)

type Credintial struct {
	User             string
	Role             string
	Password         string
	PasswordFromFile string
}

var DefaultOptions = []Option{
	WithDirect(false),
	WithTriggers(),
}

func WithCredential(c Credintial) Option {
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
