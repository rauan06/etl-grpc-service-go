package domain

import "errors"

const (
	StatusNotStarted = iota
	StatusRunning
	StatusShutdown
)

var ErrParseInt64 = errors.New("couldn't parse int64 to string")

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
