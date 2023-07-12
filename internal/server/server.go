package server

import (
	"context"
	"log"
	"net"
	pb "passwords/internal/adapters/grpc/proto"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/pkg/common/logging"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	cfg             *configuration.ServerConfiguration
	grpcService     pb.ServiceServer
	grpcServer      *grpc.Server
	AuthorizedUsers *sync.Map
	UserOperations  *sync.Map
}

func NewServer(
	cfg *configuration.ServerConfiguration,
	grpcService pb.ServiceServer,
	AuthorizedUsers *sync.Map,
	UserOperations *sync.Map) *Server {
	return &Server{
		cfg:             cfg,
		grpcService:     grpcService,
		AuthorizedUsers: AuthorizedUsers,
		UserOperations:  UserOperations,
	}
}

func (s *Server) RunListener() {
	//check necessary data
	if s.cfg.Port == "" || s.grpcService == nil {
		return
	}
	//listener
	listener, err := net.Listen("tcp", s.cfg.Port)
	logging.LogFatal(err)
	// create grpc
	s.grpcServer = grpc.NewServer()
	// register service
	pb.RegisterServiceServer(s.grpcServer, s.grpcService)
	log.Println("Start gRPC listener")

	//start grpc
	err = s.grpcServer.Serve(listener)
	logging.LogFatal(err)

	log.Println("listener stopped")
}

func (s *Server) RunUserKicker(ctx context.Context) {
	//Ticker for run user check every s.cfg.AuthTimeout seconds
	ticker := time.NewTicker(time.Duration(int64(s.cfg.AuthTimeout)) * time.Second)

	for {
		select {
		case <-ticker.C:
			{
				var kickedUsers []string
				//check users last operation time to disauthorize silent users
				s.UserOperations.Range(func(key any, val any) bool {
					//Time of an authorized user without any operation
					lastOperationTime, ok := val.(time.Time)
					if !ok {
						log.Println("unable to get time.Time type of user's val")
						return false
					}
					userBeingDurationWithoutOp := time.Now().Sub(lastOperationTime)
					if int64(userBeingDurationWithoutOp) > int64(s.cfg.AuthTimeout) {
						userName, ok := key.(string)
						if !ok {
							log.Println("unable to get string type of user key")
							return false
						}
						kickedUsers = append(kickedUsers, userName)
					}
					return true
				})
				//delete silent users from list
				for _, user := range kickedUsers {
					//Delete user from authorized
					s.AuthorizedUsers.Delete(user)
					//Delete user last operation
					s.UserOperations.Delete(user)
					log.Printf("User %v kicked the server by timeout", user)
				}
			}
		case <-ctx.Done():
			//leave loop in case main the parent context is done
			log.Println("User kicker stopped")
			break
		}
	}

}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	log.Println("Stop grpc server")
}
