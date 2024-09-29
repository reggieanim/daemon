package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogsHandler(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "stats.log")
	err := os.WriteFile(tmpFile, []byte("test log data"), 0644)
	assert.NoError(t, err)

	app := &App{
		logDir: tmpDir,
		logger: log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
	}

	t.Run("FileExists", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/logs", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.logsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, "test log data", rr.Body.String())
	})

	t.Run("FileDoesNotExist", func(t *testing.T) {
		os.Remove(tmpFile)

		req, err := http.NewRequest(http.MethodGet, "/logs", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.logsHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Failed to open log file\n", rr.Body.String())
	})
}
