package logger

import (
	"log"
	"os"
	"time"
)

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
	entry := LogEntry{
		Timestamp: time.Now(),
		NodeAddr:  nodeAddr,
		NodePort:  nodePort,
		Level:     level,
		Message:   msg,
		Type:      "system",
	}
	logToFile(entry, "system.log")
}
