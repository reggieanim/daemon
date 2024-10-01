package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func (a *serverApplication) logsHandler(w http.ResponseWriter, r *http.Request) {
	if a.logDir == "" {
		a.logDir = "."
	}

	logFilePath := filepath.Join(a.logDir, "stats.log")
	file, err := os.Open(logFilePath)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		a.logger.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	stats, err := os.ReadFile(logFilePath)
	if err != nil {
		http.Error(w, "Failed to read log file", http.StatusInternalServerError)
		a.logger.Printf("Error reading log file: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(stats)
}
