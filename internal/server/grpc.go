package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	"github.com/dkmelnik/metrics/internal/logger"
)

type GRPCServer struct {
	addr string
	app  *grpc.Server
}

func recoveryHandler(p interface{}) error {
	logger.Log.Error("server:GRPCServer", "Panic occurred:", p)
	return status.Errorf(codes.Internal, "Internal server error")
}

func NewGRPCServer(addr string) (*GRPCServer, error) {
	grpcServer := grpc.NewServer(

		grpc.ChainStreamInterceptor(),

		grpc.ChainUnaryInterceptor(),
	)

	return &GRPCServer{
		addr, grpcServer,
	}, nil
}

func (s *GRPCServer) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	go func() {
		<-sig
		logger.Log.Info("server", "stop serve addr", s.addr)
		s.app.Stop()
	}()
	logger.Log.Info("server", "serve addr", s.addr)
	return s.app.Serve(listener)
}

func (s *GRPCServer) GetApp() *grpc.Server {
	return s.app
}
