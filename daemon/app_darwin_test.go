package main

import (
	"os"
	"os/exec"
	"testing"
)

// Mocking the exec.Command function
func mockExecCommand(command string, args ...string) *exec.Cmd {
	// Simulate a successful command execution
	cs := []string{"-test.run=TestMockHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestApp_StartService(t *testing.T) {
	app := NewApp()

	msg, err := app.StartService()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedMsg := "Service started"
	if msg != expectedMsg {
		t.Fatalf("Expected %v, got %v", expectedMsg, msg)
	}

	if !app.workerRunning {
		t.Fatalf("Expected workerRunning to be true, got false")
	}
	if !app.timerRunning {
		t.Fatalf("Expected timerRunning to be true, got false")
	}
}

func TestApp_StopService(t *testing.T) {
	app := NewApp()
	_, err := app.StartService()
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	msg, err := app.StopService()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedMsg := "Service stopped"
	if msg != expectedMsg {
		t.Fatalf("Expected %v, got %v", expectedMsg, msg)
	}

	if app.workerRunning {
		t.Fatalf("Expected workerRunning to be false, got true")
	}
	if app.timerRunning {
		t.Fatalf("Expected timerRunning to be false, got true")
	}
}
