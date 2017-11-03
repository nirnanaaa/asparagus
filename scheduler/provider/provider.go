package provider

// SourceProvider is used to get tasks
type SourceProvider interface {
	Read()
}

// ExecutionProvider is used for executing cronjobs
type ExecutionProvider interface {
	Execute()
}
