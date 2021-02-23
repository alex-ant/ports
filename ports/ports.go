package ports

import (
	"context"
	"encoding/json"
	"fmt"
)

// Server defines ports gRPC server.
type Server struct {
}

func (s *Server) mustEmbedUnimplementedPortServiceServer() {}

// Ping is used for gRPC server healthcheck.
func (s *Server) Ping(context.Context, *Empty) (*Empty, error) {
	return new(Empty), nil
}

// StorePortInfo stores incoming port info into the database.
func (s *Server) StorePortInfo(ctx context.Context, pi *PortInfo) (*Empty, error) {
	piB, piBErr := json.Marshal(*pi)
	if piBErr != nil {
		return nil, fmt.Errorf("failed to marshal pi: %v", piBErr)
	}

	fmt.Println("===>> received port info:", string(piB))

	return new(Empty), nil
}
