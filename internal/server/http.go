package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	app *http.Server
}

func NewHTTPServer(addr string, r http.Handler) *HTTPServer {
	return &HTTPServer{app: &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Second * 30,
		Handler:      r,
		WriteTimeout: time.Second * 30,
	}}
}

func (s *HTTPServer) Run() {
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

	go func() {
		err := s.app.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}
