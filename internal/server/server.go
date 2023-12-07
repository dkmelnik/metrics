package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	app *http.Server
}

func NewServer(addr string) *Server {
	return &Server{app: &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}}
}

func (s *Server) Run() error {
	go s.shutdown()

	s.configureRouter()

	return s.app.ListenAndServe()
}

// shutdown listens for signals and stops the server if they arrive
func (s *Server) shutdown() error {
	quit := make(chan os.Signal, 1)

	// kill -15 <PID>, Ctrl-Z
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGTSTP)

	<-quit

	ctx, clFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer clFunc()

	return s.app.Shutdown(ctx)
}
