package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
		Addr:  addr,
		Store: kvStore,
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

		go s.handleClient(conn)
	}

	return nil
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection established from: %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		parts := strings.Fields(msg)
		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERR syntax: SET <key> <value>\n"))
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")

			s.Store.Set(key, value)
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERR syntax: GET <key>\n"))
				continue
			}
			key := parts[1]

			value, exists := s.Store.Get(key)
			if !exists {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(value + "\n"))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERR syntax: DEL <key>\n"))
				continue
			}
			key := parts[1]

			deleted := s.Store.Del(key)
			if deleted {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}

		default:
			conn.Write([]byte("ERR unknown command '" + cmd + "'\n"))
		}
	}
}
