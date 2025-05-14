package domain

const (
	StatusNotStarted = iota
	StatusRunning
	StatusShutdown
)

func StatusToString(status int) string {
	switch status {
	case StatusNotStarted:
		return "service not started"
	case StatusRunning:
		return "service is running..."
	case StatusShutdown:
		return "service is shutdown"
	default:
		return "unknown status"
	}
}
