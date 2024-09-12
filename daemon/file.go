package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/osquery/osquery-go"
)

type FileInfo struct {
	Path         string `json:"path"`
	ModifiedTime string `json:"mtime"`
	Size         int64  `json:"size"`
}

func (a *App) getFileModificationStats() (string, error) {
	if a.osqueryInstance == nil {
		return "", fmt.Errorf("osquery instance not initialized")
	}

	client, err := osquery.NewClient(a.osquerySocketPath, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to create osquery client: %w", err)
	}
	defer client.Close()

	query := fmt.Sprintf("SELECT path, mtime, size FROM file WHERE directory = '%s' ORDER BY mtime DESC", a.config.MonitorDirectory)
	response, err := client.Query(query)
	if err != nil {
		return "", fmt.Errorf("failed to execute osquery query: %w", err)
	}

	var files []FileInfo

	for _, r := range response.Response {
		mtimeUnix, err := strconv.ParseInt(r["mtime"], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse mtime: ", err)
			continue
		}
		mtime := time.Unix(mtimeUnix, 0).Format(time.RFC3339)
		size, err := strconv.ParseInt(r["size"], 10, 64)
		if err != nil {
			fmt.Println("Failed to parse size: ", err)
			continue
		}
		files = append(files, FileInfo{
			Path:         r["path"],
			ModifiedTime: mtime,
			Size:         size,
		})
	}
	jsonData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal file info to JSON: %w", err)
	}

	a.logger.Println("Updated files stats")
	return string(jsonData), nil
}

func (a *App) GetLatestFileModifications() string {
	stats, err := a.getFileModificationStats()
	if err != nil {
		return fmt.Sprintf("Error getting file modification stats: %v", err)
	}
	return stats
}

func (a *App) saveStatsToFile(fileStats string, systemStats string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	file, err := os.OpenFile("stats.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stats log file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("File Stats: %s\nSystem Stats: %s\n", fileStats, systemStats))
	if err != nil {
		return fmt.Errorf("failed to write stats to file: %w", err)
	}

	a.logger.Println("Successfully saved stats to file")
	return nil
}
