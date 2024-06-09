package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
)

type RaftNodeLog struct {
	*pb.RaftNodeLog
}

var (
	nodeAddr string
	nodePort string
	RaftLog  RaftNodeLog
	FileName string
)

func init() {
	RaftLog = RaftNodeLog{
		&pb.RaftNodeLog{
			Entries: make([]*pb.RaftLogEntry, 0),
		},
	}
}

func SetRaftAddrAndPort(addr, port string) {
	nodeAddr = addr
	nodePort = port
	FileName = fmt.Sprintf("%s_%s_%s", nodeAddr, nodePort, "raft_log.json")
}

// Init load
func InitLoadRaftLogs() {
	if FileName == "" {
		WriteSystemLog(ERROR, fmt.Sprintln("Raft address & port not specified for logs, fail to load entries"), nodeAddr, nodePort)
		return
	}
	log, err := LoadRaftLogFromFile(FileName)
	if err != nil {
		WriteSystemLog(ERROR, fmt.Sprintf("Error loading Raft log: %v", err), nodeAddr, nodePort)
		return
	}
	RaftLog.Entries = log
}

func WriteRaftLog(term uint64, command string) {
	entry := &pb.RaftLogEntry{
		Term:    term,
		Command: command,
	}

	RaftLog.Entries = append(RaftLog.Entries, entry)

	LogToFile(entry, FileName, "logs")
}

// To create a new copy log file
func (log *RaftNodeLog) SaveRaftLogToFile(filename string) error {
	if _, err := os.Stat(LogDir); os.IsNotExist(err) {
		err := os.Mkdir(LogDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating log directory: %v", err)
		}
	}

	// Fullpath with dir
	filePath := filepath.Join(LogDir, filename)

	logBytes, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("error marshaling log entry: %v", err)
	}

	err = os.WriteFile(filePath, logBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing to log file: %v", err)
	}

	return nil
}

// load log from specified file
func LoadRaftLogFromFile(filename string) ([]*pb.RaftLogEntry, error) {
	filePath := filepath.Join(LogDir, filename)

	logBytes, err := os.ReadFile(filePath)
	if err != nil {
		return []*pb.RaftLogEntry{}, fmt.Errorf("error reading log file: %v", err)
	}

	var log []*pb.RaftLogEntry
	err = json.Unmarshal(logBytes, &log)
	if err != nil {
		return []*pb.RaftLogEntry{}, fmt.Errorf("error unmarshaling log entry: %v", err)
	}

	return log, nil
}
