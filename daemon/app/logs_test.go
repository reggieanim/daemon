package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockFileHandler is a mock implementation of FileHandler for testing.
type MockFileHandler struct {
	OpenFunc     func(name string) (io.ReadCloser, error)
	ReadFileFunc func(name string) ([]byte, error)
}

func (m *MockFileHandler) Open(name string) (io.ReadCloser, error) {
	return m.OpenFunc(name)
}

func (m *MockFileHandler) ReadFile(name string) ([]byte, error) {
	return m.ReadFileFunc(name)
}

// MockReadCloser mocks an io.ReadCloser for testing.
type MockReadCloser struct {
	io.Reader
}

func (m *MockReadCloser) Close() error {
	return nil
}

func TestLogsHandler_Success(t *testing.T) {
	mockFileHandler := &MockFileHandler{
		OpenFunc: func(name string) (io.ReadCloser, error) {
			return &MockReadCloser{Reader: strings.NewReader("mock log content")}, nil
		},
		ReadFileFunc: func(name string) ([]byte, error) {
			return []byte("mock log content"), nil
		},
	}

	app := &App{
		fileHandler: mockFileHandler,
	}

	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.logsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "mock log content"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
