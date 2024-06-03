package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (level LogLevel) String() string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	NodeAddr  string    `json:"node_addr"`
	NodePort  string    `json:"node_port"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
}

func LogToFile(entry LogEntry, filename string) {
	logBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("Error marshaling log entry: %v\n", err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.Write(append(logBytes, '\n')); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
