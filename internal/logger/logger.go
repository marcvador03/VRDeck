package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	file   *os.File
	prefix string
	mu     sync.Mutex
}

var logfile string = "streamdeckVR.log"

func NewLogger(filepath, prefix string) (*Logger, error) {

	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filepath, err)
	}

	return &Logger{file: file, prefix: prefix}, nil
}

func (l *Logger) Log(message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	timestamp := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(l.file, "%s %s%s\n", timestamp, l.prefix, message)
	return err
}

func (l *Logger) Info(message string) error {
	return l.Log(fmt.Sprintf("INFO: %s", message))
}

func (l *Logger) Error(err error) error {
	return l.Log(fmt.Sprintf("ERROR: %v", err))
}

func (l *Logger) Close() error {
	return l.file.Close()
}
