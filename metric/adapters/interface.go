package adapters

// Reporter starts a new reporting agent
type Reporter interface {
	StartReporting() error
	StopReporting() error
}
