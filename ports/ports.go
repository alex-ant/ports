package ports

import (
	"context"
	"errors"
	"fmt"
)

// Storage defines storage interface.
type Storage interface {
	InsertPortInfo(pi *PortInfo) error
	FetchPortInfo() (map[string]*PortInfo, error)
}

// Server defines ports gRPC server.
type Server struct {
	storage Storage
}

// NewServer returns new Server.
func NewServer(storage Storage) (*Server, error) {
	if storage == nil {
		return nil, errors.New("empty storage interface provided")
	}

	return &Server{
		storage: storage,
	}, nil
}

func (s *Server) mustEmbedUnimplementedPortServiceServer() {}

// Ping is used for gRPC server healthcheck.
func (s *Server) Ping(context.Context, *Empty) (*Empty, error) {
	return new(Empty), nil
}

// StorePortInfo stores incoming port info into the database.
func (s *Server) StorePortInfo(ctx context.Context, pi *PortInfo) (*Empty, error) {
	insertErr := s.storage.InsertPortInfo(pi)
	if insertErr != nil {
		return nil, fmt.Errorf("failed to insert port info: %v", insertErr)
	}

	return new(Empty), nil
}

// FetchPortInfo fetches and returns all port info.
func (s *Server) FetchPortInfo(context.Context, *Empty) (*MultiplePortInfo, error) {
	pii, piiErr := s.storage.FetchPortInfo()
	if piiErr != nil {
		return nil, fmt.Errorf("failed to fetch port info from the DB: %v", piiErr)
	}

	res := make([]*PortInfo, len(pii))
	var c int
	for _, pi := range pii {
		res[c] = pi
		c++
	}

	return &MultiplePortInfo{
		Ports: res,
	}, nil
}
