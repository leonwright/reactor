package server

import (
	"context"
	"log"
	"net"

	pb "github.com/leonwright/reactor/reactorserver"
	"github.com/leonwright/reactor/utils"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedApiServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

// StartGrpcServer starts the main GRPC API.
func StartGrpcServer(cfg utils.Config) {
	deb.Infof("GRPC listening on port %d", cfg.Server.GrpcPort)
	lis, err := net.Listen("tcp", utils.GetFullHost("", cfg.Server.GrpcPort))
	if err != nil {
		utils.ProcessError(err)
	}
	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		utils.ProcessError(err)
	}
}
