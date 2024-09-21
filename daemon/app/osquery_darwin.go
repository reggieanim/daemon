package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/osquery/osquery-go"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type CommandRunner interface {
	Execute(name string, args ...string) (Command, error)
}

type Command interface {
	Start() error
	Wait() error
	StdinPipe() (io.WriteCloser, error)
	Output() *bytes.Buffer
}

type OsqueryCommandRunner struct{}

func (r *OsqueryCommandRunner) Execute(name string, args ...string) (Command, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	return &OsqueryCommand{cmd: cmd, out: &out}, nil
}

type OsqueryCommand struct {
	cmd *exec.Cmd
	out *bytes.Buffer
}

func (r *OsqueryCommand) Start() error {
	return r.cmd.Start()
}

func (r *OsqueryCommand) Wait() error {
	return r.cmd.Wait()
}

func (r *OsqueryCommand) StdinPipe() (io.WriteCloser, error) {
	return r.cmd.StdinPipe()
}

func (r *OsqueryCommand) Output() *bytes.Buffer {
	return r.out
}

// initOsquery initializes the osquery instance.
func (a *App) initOsquery() error {
	cmd, err := a.cmdRunner.Execute("osqueryi", "--nodisable_extensions")
	if err != nil {
		return fmt.Errorf("failed to execute osqueryi: %w", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin for osqueryi: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start osqueryi: %w", err)
	}

	time.Sleep(2 * time.Second)

	_, err = stdin.Write([]byte("select value from osquery_flags where name = 'extensions_socket';\n"))
	if err != nil {
		return fmt.Errorf("failed to send SQL command to osqueryi: %w", err)
	}

	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for osqueryi process: %w", err)
	}

	socketPath := a.parseSocketPath(cmd.Output().String())
	if socketPath == "" {
		return fmt.Errorf("could not retrieve socket path from osqueryi output")
	}
	a.osquerySocketPath = socketPath

	server, err := osquery.NewExtensionManagerServer("file_monitor", socketPath)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to create osquery extension: %w", err)
	}
	a.osqueryInstance = server
	go func() {
		if err := a.osqueryInstance.Run(); err != nil {
			wailsRuntime.LogErrorf(a.ctx, "osquery extension server stopped: %v", err)
		}
	}()

	return nil
}

// parseSocketPath parses the output to extract the socket path.
func (a *App) parseSocketPath(output string) string {
	re := regexp.MustCompile(`(\/[^\s\d]+(?:\.[^\s\d]+)+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// GetOsqueryStatus returns the status of osquery.
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
