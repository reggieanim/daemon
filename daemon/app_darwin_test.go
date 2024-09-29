package main

import (
	"log"
	"os"
	"testing"
)

func TestApp_StartService(t *testing.T) {
	app := NewApp()
	app.logger = log.New(os.Stdout, "TestLogger: ", log.LstdFlags)

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
	app.logger = log.New(os.Stdout, "TestLogger: ", log.LstdFlags)
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
