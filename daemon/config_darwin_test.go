package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func TestCreateDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yaml")
	monitorDir := filepath.Join(tmpDir, "monitor")

	if err := os.MkdirAll(monitorDir, os.ModePerm); err != nil {
		t.Fatalf("failed to create test monitor directory: %v", err)
	}

	app := &App{}
	err := app.createDefaultConfig(configPath, monitorDir)
	if err != nil {
		t.Fatalf("createDefaultConfig returned an error: %v", err)
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read generated config file: %v", err)
	}

	monitorDir = filepath.Clean(monitorDir)
	monitorDir = strings.ReplaceAll(monitorDir, `\`, `\\`)

	expectedConfig := `
monitor_directory: "` + monitorDir + `"
check_frequency: 60
api_endpoint: "https://eo13t4hn4shbd6x.m.pipedream.net"
`

	if strings.TrimSpace(string(configData)) != strings.TrimSpace(expectedConfig) {
		t.Errorf("config content mismatch: got %v, want %v", string(configData), expectedConfig)
	}
}

type MockDialog struct {
	mock.Mock
}

func (m *MockDialog) SaveFileDialog(ctx context.Context, options runtime.SaveDialogOptions) (string, error) {
	args := m.Called(ctx, options)
	return args.String(0), args.Error(1)
}

func (m *MockDialog) OpenDirectoryDialog(ctx context.Context, options runtime.OpenDialogOptions) (string, error) {
	args := m.Called(ctx, options)
	return args.String(0), args.Error(1)
}

func TestLoadConfig(t *testing.T) {
	mockDialogService := new(MockDialog)
	ctx := context.Background()

	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yaml")
	monitorDir := filepath.Join(tmpDir, "monitor")

	mockDialogService.On("SaveFileDialog", ctx, mock.Anything).Return(configPath, nil)
	mockDialogService.On("OpenDirectoryDialog", ctx, mock.Anything).Return(monitorDir, nil)

	err := os.MkdirAll(monitorDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create test monitor directory: %v", err)
	}

	app := &App{
		dialog: mockDialogService,
		logger: log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
	}
	err = app.loadConfig(ctx)

	assert.NoError(t, err)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created")
	}

	configData, err := os.ReadFile(configPath)
	assert.NoError(t, err)

	monitorDir = filepath.Clean(monitorDir)
	monitorDir = strings.ReplaceAll(monitorDir, `\`, `\\`)

	expectedConfig := `
monitor_directory: "` + monitorDir + `"
check_frequency: 60
api_endpoint: "https://eo13t4hn4shbd6x.m.pipedream.net"
`
	assert.Equal(t, strings.TrimSpace(expectedConfig), strings.TrimSpace(string(configData)))

}
