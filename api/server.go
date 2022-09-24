package api

import (
	"embed"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/config"
)

var TestDir embed.FS

var (
	NASA_APOD_API_URL string = "https://api.nasa.gov/planetary/apod"
	// UserIPLimiter stores IPs and boolean flags indicating whether the client IP is already in use/connected.
	UserIPLimiter = make(map[string]bool)
)

// Server serves HTTP requests for our url-collector service.
type Server struct {
	config config.Config
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Config) (*Server, error) {
	server := &Server{config: config}

	server.setupRouter()

	return server, nil
}

// setupRouter set up the router for HTTP requests.
func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("pictures", server.listPicturesURL)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(fmt.Sprintf("0.0.0.0:%s", address))
}
