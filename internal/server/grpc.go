package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/dkmelnik/metrics/internal/logger"
	grpc2 "github.com/dkmelnik/metrics/internal/middlewares/grpc"
)

type GRPCServer struct {
	addr string
	app  *grpc.Server
}

func NewGRPCServer(addr string, subnet, pkey string, signer grpc2.Signer) (*GRPCServer, error) {
	mm, err := grpc2.NewMiddlewareManager(subnet, pkey, signer)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(

		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mm.Recovery(),
				mm.TrustedSubnet(),
				mm.Decryption(),
				mm.Compress(),
			),
		),
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
		logger.Log.Info("grpc server", "stop serve addr", s.addr)
		s.app.Stop()
	}()
	logger.Log.Info("grpc server", "serve addr", s.addr)

	go func() {
		err := s.app.Serve(listener)
		if err != nil {
			logger.Log.Error("error starting grpc server", "serve addr", s.addr)
			return
		}
	}()

	return nil
}

func (s *GRPCServer) GetApp() *grpc.Server {
	return s.app
}
