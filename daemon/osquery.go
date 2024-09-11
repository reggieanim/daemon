package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/osquery/osquery-go"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) initOsquery() error {
	socketPath, err := a.discoverOsquerySocket()
	if err != nil {
		return fmt.Errorf("failed to discover osquery socket: %w", err)
	}

	a.osquerySocketPath = socketPath

	server, err := osquery.NewExtensionManagerServer("file_monitor", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create osquery extension: %w", err)
	}

	a.osqueryInstance = server

	go func() {
		if err := server.Run(); err != nil {
			wailsRuntime.LogErrorf(a.ctx, "osquery extension server stopped: %v", err)
		}
	}()

	return nil
}

func (a *App) discoverOsquerySocket() (string, error) {
	// Check environment variable first
	if envSocket := os.Getenv("OSQUERY_SOCKET"); envSocket != "" {
		return envSocket, nil
	}

	// Get current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	// Check common locations, including user-specific ones
	commonPaths := []string{
		filepath.Join(currentUser.HomeDir, ".osquery", "shell.em"),
		"/var/osquery/osquery.sock",
		"/var/run/osquery/osquery.sock",
		filepath.Join(os.TempDir(), "osquery.sock"),
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not discover osquery socket")
}

func (a *App) GetOsqueryStatus() string {
	if a.osquerySocketPath == "" {
		return "Osquery socket path not set. Has initOsquery been called?"
	}

	if _, err := os.Stat(a.osquerySocketPath); os.IsNotExist(err) {
		return fmt.Sprintf("Osquery socket not found at %s. Is osqueryd running?", a.osquerySocketPath)
	}

	if a.osqueryInstance == nil {
		return "Osquery socket found, but extension is not initialized"
	}

	return fmt.Sprintf("Osquery is running and extension is initialized. Socket: %s", a.osquerySocketPath)
}
