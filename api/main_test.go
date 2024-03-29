package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pawlobanano/NASA-Astronomy-Picture-of-the-Day-API/config"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *Server {
	config := config.Config{
		NASAAPIKey:         "DEMO_KEY",
		ConcurrentRequests: 5,
		ServerPort:         "8080",
	}

	server, err := NewServer(config)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
