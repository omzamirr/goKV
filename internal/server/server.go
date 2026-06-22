package server

import (
	"net"
	"sync"

	"github.com/omzamirr/internal/store"
)

type Server struct {
	Addr  string
	Store *store.Store
	Mu    sync.Mutex
}

func New(addr string, kvStore *store.Store) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	return nil
}
