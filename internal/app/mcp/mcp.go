package mcp

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpServer struct {
	server *mcp.Server
	host   string
}

func NewMcpServer(host string) *McpServer {
	return &McpServer{
		server: mcp.NewServer(&mcp.Implementation{
			Name:    "time-server",
			Version: "1.0.0",
		}, nil),
		host: host,
	}
}

func (s *McpServer) Start() error {
	// Add the cityTime tool.

	// Create the streamable HTTP handler.
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return s.server
	}, nil)

	s.Route()

	// Start the HTTP server in a goroutine.
	go func() {
		if err := http.ListenAndServe(s.host, handler); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	return nil
}

func (s *McpServer) Route() {
	// mcp.AddTool(s.server, &mcp.Tool{
	// 	Name:        "cityTime",
	// 	Description: "Get the current time in NYC, San Francisco, or Boston",
	// }, getTime)
}
