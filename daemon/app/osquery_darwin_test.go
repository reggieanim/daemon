package app

import (
	"bytes"
	"io"
	"testing"
)

// Mock implementation of io.WriteCloser
type MockWriteCloser struct {
	*bytes.Buffer
}

func (m *MockWriteCloser) Close() error {
	return nil
}

// Mock implementation of Command interface
type MockCommand struct {
	output *bytes.Buffer
}

func (m *MockCommand) Start() error {
	return nil
}

func (m *MockCommand) Wait() error {
	return nil
}

func (m *MockCommand) StdinPipe() (io.WriteCloser, error) {
	return &MockWriteCloser{bytes.NewBuffer(nil)}, nil
}

// MockExtensionManagerServer implements the OsqueryManager interface for testing.
type MockExtensionManagerServer struct{}

func (m *MockExtensionManagerServer) Run() error {
	return nil
}
func (m *MockCommand) Output() *bytes.Buffer {
	return m.output
}

// Mock implementation of CommandRunner interface
type MockCommandRunner struct{}

func (m *MockCommandRunner) Execute(name string, args ...string) (Command, error) {
	mockCmd := &MockCommand{
		output: bytes.NewBufferString("/Users/animr/.osquery/shell.em"),
	}
	return mockCmd, nil
}

// Test function for initOsquery
func TestInitOsquery(t *testing.T) {
	// Create an instance of App with the mock CommandRunner
	app := NewApp()
	app.cmdRunner = &MockCommandRunner{}

	// Assign the mock OsqueryManager to osqueryInstance
	app.osqueryInstance = &MockExtensionManagerServer{}

	// Run initOsquery
	err := app.initOsquery()

	// Assert that no error occurred
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Assert that the socket path was correctly set
	expectedSocketPath := "/Users/animr/.osquery/shell.em"
	if app.osquerySocketPath != expectedSocketPath {
		t.Errorf("Expected socket path %s, but got %s", expectedSocketPath, app.osquerySocketPath)
	}
}
