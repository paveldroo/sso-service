package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/paveldroo/sso-service/internal/config"
	pb "github.com/paveldroo/sso-service/protos/sso"
	"google.golang.org/grpc"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
}

func New() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &AuthServer{})
	return s
}

func (a AuthServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{Message: "OK"}, nil
}

func StartServer(cfg config.Config, s *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err = s.Serve(lis); err != nil {
		slog.Error("running server: %w", err)
	}
}
