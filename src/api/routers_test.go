package api_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"mistapi/src/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	t.Parallel()

	// ARANGE
	mockClient := new(MockClient)
	MockGrpcClient(t, mockClient)

	// Set test env port
	os.Setenv("APP_PORT", "8081")

	// ACT
	go func() {
		api.StartService() // This will block so run in goroutine
	}()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	// ASSERT
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
