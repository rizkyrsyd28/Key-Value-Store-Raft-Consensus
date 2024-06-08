package test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/logger"
)

func TestRaftLog(t *testing.T) {
	nodeAddr := "127.0.0.1"
	nodePort := "8080"

	tempFileName := fmt.Sprintf("%s_%s_%s", nodeAddr, nodePort, "raft_log.json")

	// Create file somewhere

	// Log entries
	logger.SetRaftAddrAndPort(nodeAddr, nodePort)
	// logger.InitLoadRaftLogs()
	logger.WriteRaftLog(69, "PING")
	logger.WriteRaftLog(420, "SET abob bibob")

	// Save the logs
	// err := logger.RaftLog.SaveRaftLogToFile(tempFileName)
	// if err != nil {
	// 	t.Fatalf("Failed to save Raft log: %v", err)
	// }

	logEntries, err := logger.LoadRaftLogFromFile(tempFileName)
	if err != nil {
		t.Fatalf("Failed to load Raft log: %v", err)
	}

	expected := logger.RaftNodeLog{
		{
			Term:    69,
			Command: "PING",
		},
		{
			Term:    420,
			Command: "SET abob bibob",
		},
	}

	if !equalLogs(expected, logEntries) {
		t.Errorf("Raft log entries mismatch\nExpected: %v\nGot: %v", expected, logEntries)
	}

	// Clean up
	err = os.Remove(logger.LogDir+"/"+tempFileName)
	if err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

func equalLogs(expected, actual logger.RaftNodeLog) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		if expected[i].Term != actual[i].Term ||
			expected[i].Command != actual[i].Command {
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
	// Data
	level := logger.INFO
	msg := "Test message"
	nodeAddr := "127.0.0.1"
	nodePort := "8080"

	logFilePath := fmt.Sprintf(logger.LogDir+"/%s_%s_system_log.json", nodeAddr, nodePort)

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
		var entry []logger.SystemLogEntry
		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			t.Fatalf("failed to unmarshal log entry: %s", err)
		}
		if entry[0].Level == level && entry[0].Message == msg && entry[0].NodeAddr == nodeAddr && entry[0].NodePort == nodePort {
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
