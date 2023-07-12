package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "passwords/internal/adapters/grpc/proto"
	"passwords/internal/pkg/common/logging"
)

type GRPCClient struct {
	Client pb.ServiceClient
	conn   *grpc.ClientConn
}

func (g GRPCClient) Close() {
	g.conn.Close()
}

func NewGRPCClient(address string) *GRPCClient {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	logging.LogFatal(err)

	return &GRPCClient{
		Client: pb.NewServiceClient(conn),
		conn:   conn}
}
