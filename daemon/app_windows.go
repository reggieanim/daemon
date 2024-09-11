package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/osquery/osquery-go"
)

type App struct {
	ctx               context.Context
	config            Config
	logger            *log.Logger
	osquerySocketPath string
	osqueryInstance   *osquery.ExtensionManagerServer
	workerQueue       chan string
	timerLogs         []string
	mutex             sync.Mutex
}

func NewApp() *App {
	return &App{
		workerQueue: make(chan string, 100),
		timerLogs:   []string{},
		logger:      log.New(os.Stdout, "AppLogger: ", log.LstdFlags),
	}
}

func (a *App) startup(ctx context.Context) {
	err := a.loadConfig(ctx)
	if err != nil {
		a.logger.Println("Could not load config", err)
	}
	a.ctx = ctx
	err = a.initOsquery()
	if err != nil {
		a.logger.Println("Could not connect to osquery")
	}
	go a.workerThread()
	go a.timerThread()
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Windows-friendly workerThread that uses "cmd /C" for command execution
func (a *App) workerThread() {
	for cmdStr := range a.workerQueue {
		a.logger.Printf("Executing command: %s", cmdStr)

		// On Windows, use "cmd /C" to execute shell commands
		cmd := exec.Command("cmd", "/C", cmdStr)

		// Capture the output (stdout and stderr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			a.logger.Printf("Error executing command: %v, output: %s", err, output)
			continue
		}

		a.logger.Printf("Command output: %s", output)
	}
}

func (a *App) timerThread() {
	ticker := time.NewTicker(time.Duration(a.config.CheckFrequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats, err := a.getFileModificationStats()
		if err != nil {
			a.logger.Printf("Error getting file modification stats: %v", err)
			continue
		}

		systemStats, err := a.getSystemMonitoringData()
		fmt.Println(systemStats)
		if err != nil {
			a.logger.Printf("Error getting system monitoring data: %v", err)
			continue
		}
		a.mutex.Lock()
		a.timerLogs = append(a.timerLogs, stats)
		a.mutex.Unlock()

		if err := a.sendStatsToAPI(stats, systemStats); err != nil {
			a.logger.Printf("Error sending stats to API: %v", err)
		}
	}
}
