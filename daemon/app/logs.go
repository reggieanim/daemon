package app

import (
	"net/http"
)

func (a *App) logsHandler(w http.ResponseWriter, r *http.Request) {
	file, err := a.fileHandler.Open("stats.log")
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		a.logger.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	stats, err := a.fileHandler.ReadFile("stats.log")
	if err != nil {
		http.Error(w, "Failed to read log file", http.StatusInternalServerError)
		a.logger.Printf("Error reading log file: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(stats)
}
