package fb_nbackup

type Option func(m *Manager)

func WithCredintial() Option {
	return func(m *Manager) {
	}
}

func WithDirect(use bool) Option {
	return func(m *Manager) {
	}
}

func WithoutTriggers() Option {
	return func(m *Manager) {
	}
}

func WithTriggers() Option {
	return func(m *Manager) {
	}
}

func DefaultOptions() []Option {
	return []Option{
		WithDirect(false),
		WithTriggers(),
	}
}
