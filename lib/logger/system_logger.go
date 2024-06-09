package logger

import (
	"fmt"
	"log"
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

type SystemLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	NodeAddr  string    `json:"node_addr"`
	NodePort  string    `json:"node_port"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
}

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
)

func init() {
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	FatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime)
}

func WriteSystemLog(level LogLevel, msg string, nodeAddr string, nodePort string) {
	entry := SystemLogEntry{
		Timestamp: time.Now(),
		NodeAddr:  nodeAddr,
		NodePort:  nodePort,
		Level:     level,
		Message:   msg,
		Type:      "system",
	}
	LogToFile(entry, fmt.Sprintf("%s_%s_%s", nodeAddr, nodePort, "system_log.json"), LogDir)
}

func SetSystemLogFlags(flags int) {
	DebugLogger.SetFlags(flags)
	InfoLogger.SetFlags(flags)
	WarningLogger.SetFlags(flags)
	ErrorLogger.SetFlags(flags)
	FatalLogger.SetFlags(flags)
}
