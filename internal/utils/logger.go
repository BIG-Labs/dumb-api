package utils

import "log"

// Logger interface (if needed)
type Logger interface {
	Info(msg string)
	Error(msg string)
}

type StdLogger struct{}

func (l *StdLogger) Info(msg string) {
	log.Println("[INFO]", msg)
}

func (l *StdLogger) Error(msg string) {
	log.Println("[ERROR]", msg)
}

// TODO: Implement other loggers or custom formatting if needed.
