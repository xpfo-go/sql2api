package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/xpfo-go/logs"
)

const (
	defaultGraceTimeout = 30 * time.Second
	defaultIdleTimeout  = 180 * time.Second
	defaultReadTimeout  = 60 * time.Second
	defaultWriteTimeout = 60 * time.Second
)

var server *Server

type Config struct {
	Host string
	Port int
}

// Server ...
type Server struct {
	addr     string
	server   *http.Server
	stopChan chan struct{}
}

// NewServer ...
func NewServer(cfg Config) *Server {
	if server != nil {
		return server
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logs.Infof("the server addr: %s", addr)

	// parse the timeouts
	readTimeout := defaultReadTimeout
	writeTimeout := defaultWriteTimeout
	idleTimeout := defaultIdleTimeout

	logs.Infof("the server timeout settings: read_timeout=%s, write_timeout=%s, idle_timeout=%s",
		readTimeout, writeTimeout, idleTimeout)

	router := GetRouter()
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	server = &Server{
		addr:     addr,
		server:   httpServer,
		stopChan: make(chan struct{}, 1),
	}
	return server
}

// Run ...
func (s *Server) Run(ctx context.Context) {

	go func() {
		<-ctx.Done()
		logs.Info("I have to go...")
		logs.Info("Stopping server gracefully")
		s.Stop()
	}()

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	s.Wait()
	logs.Info("Shutting down")
}

// Stop ...
func (s *Server) Stop() {
	defer logs.Info("Server stopped")

	// default graceTimeOut is 60 seconds
	graceTimeout := defaultGraceTimeout

	ctx, cancel := context.WithTimeout(context.Background(), graceTimeout)
	defer cancel()
	logs.Infof("Waiting %s seconds before killing connections...", graceTimeout)

	// disable keep-alive connections
	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		logs.Error(err.Error())
		_ = s.server.Close()
	}

	s.stopChan <- struct{}{}
}

// Wait blocks until server is shut down.
func (s *Server) Wait() {
	<-s.stopChan
}
