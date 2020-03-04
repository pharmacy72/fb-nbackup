package fb_nbackup

type Manager struct {
}

func NewManager(opts ...Option) *Manager {
	manager := &Manager{}
	for _, option := range opts {
		option(manager)
	}
	return manager
}
