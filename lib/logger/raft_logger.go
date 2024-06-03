package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RaftNodeLog []LogEntry

var (
	nodeAddr string
	nodePort string
	RaftLog  RaftNodeLog
	FileName string
)

func SetRaftFileName(filename string) {
	FileName = filename
}

func SetRaftAddrAndPort(addr, port string) {
	nodeAddr = addr
	nodePort = port
}

// Init load
func InitLoadRaftLogs() {
	if FileName == "" {
		WriteSystemLog(ERROR, fmt.Sprintln("Raft log file name not set, fail to load entries"), nodeAddr, nodePort)
		return
	}
	log, err := LoadRaftLogFromFile(FileName)
	if err != nil {
		WriteSystemLog(ERROR, fmt.Sprintf("Error loading Raft log: %v", err), nodeAddr, nodePort)
		return
	}
	RaftLog = log
}

func WriteRaftLog(level LogLevel, msg string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		NodeAddr:  nodeAddr,
		NodePort:  nodePort,
		Level:     level,
		Message:   msg,
		Type:      "raft",
	}
	RaftLog = append(RaftLog, entry)
	LogToFile(entry, FileName)
}

// To create a new copy log file
func (log *RaftNodeLog) SaveRaftLogToFile(filename string) error {
	logBytes, err := json.Marshal(log)
	if err != nil {
		return err
	}

	// Full Dir Path
	err = os.WriteFile(filename, logBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// load log from specified file
func LoadRaftLogFromFile(filename string) (RaftNodeLog, error) {
	logBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var log RaftNodeLog
	err = json.Unmarshal(logBytes, &log)
	if err != nil {
		return nil, err
	}

	return log, nil
}
