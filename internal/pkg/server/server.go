package server

import (
	"log"
	"net"
	pb "passwords/internal/adapters/grpc/proto"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/pkg/common/logging"

	"google.golang.org/grpc"
)

type Server struct {
	cfg         *configuration.ServerConfiguration
	grpcService pb.ServiceServer
	grpcServer  *grpc.Server
}

func NewServer(cfg *configuration.ServerConfiguration, grpcService pb.ServiceServer) *Server {
	return &Server{
		cfg:         cfg,
		grpcService: grpcService,
	}
}

func (s Server) Run() {
	//check necessary data
	if s.cfg.Port == "" || s.grpcService == nil {
		return
	}
	//listener
	listener, err := net.Listen("tcp", ":"+s.cfg.Port)
	logging.LogFatal(err)
	// create grpc
	s.grpcServer = grpc.NewServer()
	// register service
	pb.RegisterServiceServer(s.grpcServer, s.grpcService)
	log.Println("Start gRPC server")
	//start
	err = s.grpcServer.Serve(listener)
	logging.LogFatal(err)
}

func (s Server) Stop() {
	s.grpcServer.GracefulStop()
	log.Println("Stop http server")
}
