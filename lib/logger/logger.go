package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
    LogDir = "logs"
)

func LogToFile(entry interface{}, filename string, logDir string) {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			fmt.Printf("Error creating log directory: %v\n", err)
			return
		}
	}

	filePath := filepath.Join(logDir, filename)

	var logs []interface{}
	
	// Read existing log file if it exists
	if _, err := os.Stat(filePath); err == nil {
		logFile, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading log file: %v\n", err)
			return
		}

		err = json.Unmarshal(logFile, &logs)
		if err != nil {
			fmt.Printf("Error unmarshaling log file: %v\n", err)
			return
		}
	}

	// Append new entry
	logs = append(logs, entry)

	// Marshal the logs slice back to JSON
	logBytes, err := json.Marshal(logs)
	if err != nil {
		fmt.Printf("Error marshaling log entries: %v\n", err)
		return
	}

	// Write the updated log entries back to the file
	err = os.WriteFile(filePath, logBytes, 0644)
	if err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
