package server

import (
	"bufio"
	"fmt"
	"log"
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection established from: %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("Client sent: %s\n", msg)
		conn.Write([]byte("Received: " + msg + "\n"))
	}

}
