package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type ServerHTTP struct {
	serverRunning bool
	sync.Mutex
	httpServer *http.Server
	port       string
	startError chan error
}

// SetPort implements Server.
func (s *ServerHTTP) SetPort(port string) {
	s.Lock()
	defer s.Unlock()
	s.port = port
}

func (s *ServerHTTP) GetStartError() <-chan error {
	s.Lock()
	defer s.Unlock()

	return s.startError
}

// IsValid implements Server.
func (s *ServerHTTP) IsValid() error {
	if err := s.isPortValid(); err != nil {
		return err
	}

	if err := s.isPortAvailability(); err != nil {
		return err
	}
	return nil
}

// StartAsync implements Server.
func (s *ServerHTTP) StartAsync() {
	s.Lock()
	defer s.Unlock()
	if s.serverRunning {
		return // Already running
	}

	addr := "localhost:" + s.port
	s.httpServer = &http.Server{Addr: addr}
	s.serverRunning = true
	s.startError = nil
	s.startError = make(chan error, 1)

	go func() {
		log.Printf("Servidor HTTP iniciado em http://%s\n", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.serverRunning = false
			s.startError <- fmt.Errorf("Erro ao iniciar servidor HTTP: %v", err)
		}
	}()
	s.serverRunning = true
	return
}

// Stop implements Server.
func (s *ServerHTTP) Stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.serverRunning {
		return nil
	}

	log.Println("Parando servidor HTTP...")
	if s.httpServer != nil {
		err := s.httpServer.Close()
		if err != nil {
			return fmt.Errorf("Erro ao parar servidor HTTP: %v", err)
		}
	}

	s.serverRunning = false
	return nil
}

func NewServerHTTP(port string) Server {
	return &ServerHTTP{
		port: port,
	}
}

func (s *ServerHTTP) isPortValid() error {
	port, err := strconv.Atoi(s.port)
	if err != nil {
		return fmt.Errorf("Porta inválida, aceita somente números inteiros") // Não é um número válido
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("Porta fora do intervalo aceito de %d a %d", 0, 65535)
	}
	return nil
}

func (s *ServerHTTP) isPortAvailability() error {
	addr := net.JoinHostPort("localhost", s.port)
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		conn.Close()
		return fmt.Errorf("Porta %s já está em uso", s.port)
	}
	return nil
}
