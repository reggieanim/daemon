package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendStatsToAPI_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		var payload map[string]string
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		require.Equal(t, "mock_file_stats", payload["file_stats"])
		require.Equal(t, "mock_system_monitor", payload["system_monitor"])

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	app := &App{
		config: Config{
			APIEndpoint: mockServer.URL,
		},
		logger: log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
	}

	err := app.sendStatsToAPI("mock_file_stats", "mock_system_monitor")

	assert.NoError(t, err)
}

func TestSendStatsToAPI_Failure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	app := &App{
		config: Config{
			APIEndpoint: mockServer.URL, // Use the mock server's URL as the API endpoint
		},
		logger: log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
	}

	err := app.sendStatsToAPI("mock_file_stats", "mock_system_monitor")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API responded with status: 500 Internal Server Error")
}
