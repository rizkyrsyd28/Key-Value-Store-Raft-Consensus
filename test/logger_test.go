package test

import (
	"encoding/json"
	"os"
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
