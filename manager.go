package fb_nbackup

const (
	defaultCommand = "nbackup"
)

type Manager struct {
	command           string
	decompressCommand string
	direct            bool
	noDBTriggers      bool
	credintial        Credintial
}

func NewManager(opts ...Option) *Manager {
	manager := &Manager{}
	for _, option := range opts {
		option(manager)
	}
	return manager
}

func (m *Manager) exec(args ...string) ([]byte, error) {
	return nil, nil
}

func (m *Manager) Version() (string, error) {
	data, err := m.exec("-Z")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
