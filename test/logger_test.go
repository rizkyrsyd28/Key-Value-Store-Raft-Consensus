package test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
)

func TestRaftLog(t *testing.T) {
	nodeAddr := "127.0.0.1"
	nodePort := "8080"

	// Create file somewhere
	tempFile, err := os.CreateTemp("", "raft_log_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()

	// Log entries
	logger.SetRaftFileName(tempFile.Name())
	logger.SetRaftAddrAndPort(nodeAddr, nodePort)
	// logger.InitLoadRaftLogs()
	logger.WriteRaftLog(logger.INFO, "Test log entry 1")
	logger.WriteRaftLog(logger.WARNING, "Test log entry 2")

	// Save the logs
	err = logger.RaftLog.SaveRaftLogToFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to save Raft log: %v", err)
	}

	// Read the logs from file
	logBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read Raft log file: %v", err)
	}

	// Unmarsh
	var logEntries logger.RaftNodeLog
	err = json.Unmarshal(logBytes, &logEntries)
	if err != nil {
		t.Fatalf("Failed to unmarshal Raft log: %v", err)
	}

	expected := logger.RaftNodeLog{
		{
			Timestamp: logEntries[0].Timestamp,
			NodeAddr:  nodeAddr,
			NodePort:  nodePort,
			Level:     logger.INFO,
			Message:   "Test log entry 1",
			Type:      "raft",
		},
		{
			Timestamp: logEntries[1].Timestamp,
			NodeAddr:  nodeAddr,
			NodePort:  nodePort,
			Level:     logger.WARNING,
			Message:   "Test log entry 2",
			Type:      "raft",
		},
	}

	if !equalLogs(expected, logEntries) {
		t.Errorf("Raft log entries mismatch\nExpected: %v\nGot: %v", expected, logEntries)
	}
}

func equalLogs(expected, actual logger.RaftNodeLog) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		if expected[i].Level != actual[i].Level ||
			expected[i].Message != actual[i].Message ||
			expected[i].NodeAddr != actual[i].NodeAddr ||
			expected[i].Type != actual[i].Type {
			return false
		}
	}

	return true
}

func TestSystemLog(t *testing.T) {
	// Disable date time flag
	logger.SetSystemLogFlags(0)
    var buf bytes.Buffer
    logger.DebugLogger.SetOutput(&buf)
    logger.InfoLogger.SetOutput(&buf)
    logger.WarningLogger.SetOutput(&buf)
    logger.ErrorLogger.SetOutput(&buf)
    logger.FatalLogger.SetOutput(&buf)

    // Log some messages
    logger.DebugLogger.Println("This is a debug message")
    logger.InfoLogger.Println("This is an info message")
    logger.WarningLogger.Println("This is a warning message")
    logger.ErrorLogger.Println("This is an error message")
    logger.FatalLogger.Println("This is a fatal message")

    // Check
    output := buf.String()
    expectedOutput := strings.Join([]string{
        "DEBUG: This is a debug message",
        "INFO: This is an info message",
        "WARNING: This is a warning message",
        "ERROR: This is an error message",
        "FATAL: This is a fatal message",
    }, "\n") + "\n"

    if output != expectedOutput {
        t.Errorf("System log output mismatch\nExpected:\n%s\nGot:\n%s", expectedOutput, output)
    }
}


func TestLoggerInitialization(t *testing.T) {
	if logger.DebugLogger == nil {
		t.Error("DebugLogger not initialized")
	}
	if logger.InfoLogger == nil {
		t.Error("InfoLogger not initialized")
	}
	if logger.WarningLogger == nil {
		t.Error("WarningLogger not initialized")
	}
	if logger.ErrorLogger == nil {
		t.Error("ErrorLogger not initialized")
	}
	if logger.FatalLogger == nil {
		t.Error("FatalLogger not initialized")
	}
}

func TestWriteSystemLog(t *testing.T) {
	logFilePath := "system.log"
	defer os.Remove(logFilePath) // Clean up later

	// Data
	level := logger.INFO
	msg := "Test message"
	nodeAddr := "127.0.0.1"
	nodePort := "8080"

	logger.WriteSystemLog(level, msg, nodeAddr, nodePort)

	file, err := os.Open(logFilePath)
	if err != nil {
		t.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	// Check
	found := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry logger.LogEntry
		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			t.Fatalf("failed to unmarshal log entry: %s", err)
		}
		if entry.Level == level && entry.Message == msg && entry.NodeAddr == nodeAddr && entry.NodePort == nodePort {
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	if !found {
		t.Errorf("expected log entry not found in file")
	}
}