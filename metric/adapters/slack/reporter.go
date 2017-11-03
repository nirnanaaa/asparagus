package slack

// Reporter type
type Reporter struct {
	Config Config
}

// NewReporter returns a new reporter
func NewReporter(config Config) *Reporter {
	return &Reporter{
		config,
	}
}

// StartReporting starts the recording process
func (r *Reporter) StartReporting() error {
	if !r.Config.Enabled {
		return nil
	}
	return nil
}

// StopReporting stops the recording process
func (r *Reporter) StopReporting() error {
	return nil
}
