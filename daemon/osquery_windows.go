package main

import (
	"fmt"
	"os"

	"github.com/osquery/osquery-go"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) initOsquery() error {
	// Discover the osquery named pipe on Windows
	namedPipePath, err := a.discoverOsqueryPipe()
	if err != nil {
		return fmt.Errorf("failed to discover osquery named pipe: %w", err)
	}

	a.osquerySocketPath = namedPipePath

	// Create an osquery extension manager server using the named pipe path
	server, err := osquery.NewExtensionManagerServer("file_monitor", namedPipePath)
	if err != nil {
		return fmt.Errorf("failed to create osquery extension: %w", err)
	}

	a.osqueryInstance = server

	// Run the osquery extension manager server
	go func() {
		if err := server.Run(); err != nil {
			wailsRuntime.LogErrorf(a.ctx, "osquery extension server stopped: %v", err)
		}
	}()

	return nil
}

// discoverOsqueryPipe checks for the osquery named pipe on Windows
func (a *App) discoverOsqueryPipe() (string, error) {
	// Named pipe used by osquery on Windows
	namedPipe := `\\.\pipe\osquery.em`

	// Check if the named pipe exists
	if _, err := os.Stat(namedPipe); err == nil {
		return namedPipe, nil
	}

	return "", fmt.Errorf("could not find osquery named pipe at %s", namedPipe)
}

func (a *App) GetOsqueryStatus() string {
	if a.osquerySocketPath == "" {
		return "Osquery named pipe not set. Has initOsquery been called?"
	}

	if _, err := os.Stat(a.osquerySocketPath); os.IsNotExist(err) {
		return fmt.Sprintf("Osquery named pipe not found at %s. Is osqueryd running?", a.osquerySocketPath)
	}

	if a.osqueryInstance == nil {
		return "Osquery named pipe found, but extension is not initialized"
	}

	return fmt.Sprintf("Osquery is running and extension is initialized. Pipe: %s", a.osquerySocketPath)
}
