package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/paveldroo/sso-service/internal/lib/logger/sl"
	pb "github.com/paveldroo/sso-service/protos/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	AddUser(email, password string, isAdmin bool) error
	User(email, password string) (int64, error)
	IsAdmin(email string) (bool, error)
}

type AuthServer struct {
	storage Storage
	pb.UnimplementedAuthServer
}

func New(storage Storage) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &AuthServer{storage: storage})
	return s
}

func (a AuthServer) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := a.storage.AddUser(r.Email, r.Password, false); err != nil {
		slog.Error("add user to storage: %w", sl.Err(err))
		return nil, status.Error(codes.Internal, "user register failed")
	}
	return &pb.RegisterResponse{Message: "OK"}, nil
}

func (a AuthServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	id, err := a.storage.User(r.Email, r.Password)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "incorrect email or password")
	}
	return &pb.LoginResponse{Token: fmt.Sprintf("userID: %d", id)}, nil
}

func (a AuthServer) IsAdmin(ctx context.Context, r *pb.IsAdminRequest) (*pb.IsAdminResponse, error) {
	isAdmin, err := a.storage.IsAdmin(r.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &pb.IsAdminResponse{IsAdmin: isAdmin}, nil
}
