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

func NewServer(addr string, r http.Handler) *Server {
	return &Server{app: &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second * 30,
		Handler:      r,
		WriteTimeout: time.Second * 30,
	}}
}

func (s *Server) Run() error {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	go func() {
		<-sig
		err := s.app.Shutdown(serverCtx)
		if err != nil {
			panic(err)
		}
		serverStopCtx()
	}()

	return s.app.ListenAndServe()
}
