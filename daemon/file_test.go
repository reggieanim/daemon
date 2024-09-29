package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/osquery/osquery-go/gen/osquery"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// type MockOsqueryClient struct {
// 	mock.Mock
// }

// func (m *MockOsqueryClient) Query(query string) (*osquery.ExtensionResponse, error) {
// 	args := m.Called(query)
// 	return args.Get(0).(*osquery.ExtensionResponse), args.Error(1)
// }

// func (m *MockOsqueryClient) Close() {
// }

// func TestGetFileModificationStats_Success(t *testing.T) {
// 	mockClient := new(MockOsqueryClient)

// 	mockResponse := osquery.ExtensionResponse{
// 		Response: []map[string]string{
// 			{"path": "/test/file1.txt", "mtime": "1632700000", "size": "1024"},
// 			{"path": "/test/file2.txt", "mtime": "1632800000", "size": "2048"},
// 		},
// 	}
// 	mockClient.On("Query", mock.Anything).Return(mockResponse, nil)
// 	app := &App{
// 		osquerySocketPath: "/tmp/osquery.sock",
// 		client:            mockClient, // Inject the mock client
// 		config:            Config{MonitorDirectory: "/test/directory"},
// 	}

// 	stats, err := app.getFileModificationStats()

// 	assert.NoError(t, err)
// 	// Expected JSON result
// 	expectedJSON := `[{
// 		"path": "/test/file1.txt",
// 		"mtime": "2021-09-26T01:06:40Z",
// 		"size": 1024
// 	}, {
// 		"path": "/test/file2.txt",
// 		"mtime": "2021-09-27T01:06:40Z",
// 		"size": 2048
// 	}]`

// 	// Assert that the JSON output matches the expected result
// 	assert.JSONEq(t, expectedJSON, stats)

// 	mockClient.AssertCalled(t, "Query", mock.Anything)
// }

// func TestGetFileModificationStats_ClientNotInitialized(t *testing.T) {
// 	app := &App{
// 		config: Config{MonitorDirectory: "/test/directory"},
// 		logger: log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
// 	}

// 	stats, err := app.getFileModificationStats()

// 	assert.Error(t, err)
// 	assert.Equal(t, "osquery instance not initialized", err.Error())

// 	assert.Empty(t, stats)
// }

// func TestGetFileModificationStats_QueryFailure(t *testing.T) {
// 	// Create a mock osquery client
// 	mockClient := new(MockOsqueryClient)

// 	// Setup the mock Query method to return an error
// 	mockClient.On("Query", mock.Anything).Return(osquery.ExtensionResponse{}, fmt.Errorf("query failed"))

// 	// Create the app instance and inject the mock client
// 	app := &App{
// 		osquerySocketPath: "/tmp/osquery.sock",
// 		client:            mockClient, // Inject the mock client
// 		config:            Config{MonitorDirectory: "/test/directory"},
// 		logger:            log.New(os.Stdout, "TestLogger: ", log.LstdFlags),
// 	}

// 	// Call the getFileModificationStats method
// 	stats, err := app.getFileModificationStats()

// 	// Assert that an error is returned
// 	assert.Error(t, err)
// 	assert.Equal(t, "failed to execute osquery query: query failed", err.Error())

// 	// Assert that stats is empty
// 	assert.Empty(t, stats)

// 	// Ensure that the Query method was called
// 	mockClient.AssertCalled(t, "Query", mock.Anything)
// }
