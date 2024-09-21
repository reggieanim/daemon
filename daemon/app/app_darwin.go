package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type OsqueryManager interface {
	Run() error
}

type CommandPayload struct {
	Command string `json:"command"`
}

type App struct {
	fileHandler       FileHandler
	cmdRunner         CommandRunner
	ctx               context.Context
	config            Config
	logger            Logger
	osquerySocketPath string
	osqueryInstance   OsqueryManager
	workerQueue       chan string
	timerLogs         []string
	mutex             sync.Mutex

	stopWorker    chan struct{}
	stopTimer     chan struct{}
	workerRunning bool
	timerRunning  bool
}

func NewApp() *App {
	return &App{
		workerQueue: make(chan string, 100),
		cmdRunner:   &OsqueryCommandRunner{},
		timerLogs:   []string{},
		logger:      log.New(os.Stdout, "AppLogger: ", log.LstdFlags),
		stopWorker:  make(chan struct{}),
		stopTimer:   make(chan struct{}),
	}
}

func (a *App) Startup(ctx context.Context) {
	err := a.loadConfig(ctx)
	if err != nil {
		a.logger.Println("Could not load config", err)
	}
	a.ctx = ctx
	err = a.initOsquery()
	if err != nil {
		a.logger.Println("Could not connect to osquery", err)
	}
}

func (a *App) workerThread() {
	a.logger.Println("Worker thread started")
	for {
		select {
		case cmdStr := <-a.workerQueue:
			a.logger.Printf("Executing command: %s", cmdStr)
			cmd := exec.Command("sh", "-c", cmdStr)
			output, err := cmd.CombinedOutput()
			if err != nil {
				a.logger.Printf("Error executing command: %v, output: %s", err, output)
				continue
			}
			a.logger.Printf("Command output: %s", output)
		case <-a.stopWorker:
			a.logger.Println("Worker thread stopped")
			return
		}
	}
}

func (a *App) StartService() (string, error) {
	if !a.workerRunning {
		go a.workerThread()
		a.workerRunning = true
	}
	if !a.timerRunning {
		go a.timerThread()
		a.timerRunning = true
	}
	return "Service started", nil
}

func (a *App) StopService() (string, error) {
	if a.workerRunning {
		a.stopWorker <- struct{}{}
		a.workerRunning = false
	}
	if a.timerRunning {
		a.stopTimer <- struct{}{}
		a.timerRunning = false
	}
	return "Service stopped", nil
}

func (a *App) timerThread() {
	var frequency int
	if a.config.CheckFrequency == 0 {
		frequency = 1
	} else {
		frequency = a.config.CheckFrequency
	}
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	defer ticker.Stop()
	a.logger.Println("Timer thread started")

	for {
		select {
		case <-ticker.C:
			stats, err := a.getFileModificationStats()
			if err != nil {
				a.logger.Printf("Error getting file modification stats: %v", err)
				continue
			}

			systemStats, err := a.getSystemMonitoringData()
			if err != nil {
				a.logger.Printf("Error getting system monitoring data: %v", err)
				continue
			}

			a.mutex.Lock()
			a.timerLogs = append(a.timerLogs, stats)
			a.mutex.Unlock()

			if err := a.saveStatsToFile(stats, systemStats); err != nil {
				a.logger.Printf("Error saving stats to file: %v", err)
			}

			if err := a.sendStatsToAPI(stats, systemStats); err != nil {
				a.logger.Printf("Error sending stats to API: %v", err)
			}
		case <-a.stopTimer:
			a.logger.Println("Timer thread stopped")
			return
		}
	}
}

func (a *App) cpuCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload CommandPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Command == "" {
		http.Error(w, "Command is required", http.StatusBadRequest)
		return
	}

	a.workerQueue <- payload.Command
	a.logger.Printf("Command enqueued: %s", payload.Command)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Command enqueued successfully"))
}
