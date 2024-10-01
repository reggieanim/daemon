package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// type MockLogger struct{}

// func TestCpuCommandHandler(t *testing.T) {
// 	workerQueue := make(chan string, 1)

// 	app := &App{
// 		workerQueue: workerQueue,
// 		logger:      log.New(os.Stdout, "AppLogger: ", log.LstdFlags),
// 	}

// 	t.Run("InvalidMethod", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodGet, "/cpu-command", nil)
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(app.cpuCommandHandler)
// 		handler.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
// 		assert.Equal(t, "Invalid request method\n", rr.Body.String())
// 	})

// 	t.Run("InvalidRequestBody", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodPost, "/cpu-command", bytes.NewBuffer([]byte("invalid body")))
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(app.cpuCommandHandler)
// 		handler.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Equal(t, "Invalid request body\n", rr.Body.String())
// 	})

// 	t.Run("MissingCommand", func(t *testing.T) {
// 		payload := CommandPayload{Command: ""}
// 		body, err := json.Marshal(payload)
// 		assert.NoError(t, err)

// 		req, err := http.NewRequest(http.MethodPost, "/cpu-command", bytes.NewBuffer(body))
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(app.cpuCommandHandler)
// 		handler.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusBadRequest, rr.Code)
// 		assert.Equal(t, "Command is required\n", rr.Body.String())
// 	})

// 	t.Run("ValidCommand", func(t *testing.T) {
// 		payload := CommandPayload{Command: "test-command"}
// 		body, err := json.Marshal(payload)
// 		assert.NoError(t, err)

// 		req, err := http.NewRequest(http.MethodPost, "/cpu-command", bytes.NewBuffer(body))
// 		assert.NoError(t, err)

// 		rr := httptest.NewRecorder()
// 		handler := http.HandlerFunc(app.cpuCommandHandler)
// 		handler.ServeHTTP(rr, req)

// 		assert.Equal(t, http.StatusOK, rr.Code)
// 		assert.Equal(t, "Command enqueued successfully", rr.Body.String())

// 		command := <-workerQueue
// 		assert.Equal(t, "test-command", command)
// 	})
// }
