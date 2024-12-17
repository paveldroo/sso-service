package authgrpc

import (
	"context"
	"errors"

	"github.com/paveldroo/sso-service/internal/service/auth"
	"github.com/paveldroo/sso-service/internal/storage"
	ssov1 "github.com/paveldroo/sso-service/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (string, error)
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Server struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(g *grpc.Server, authService Auth) {
	ssov1.RegisterAuthServer(g, &Server{auth: authService})
}

func (s Server) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if in.AppID == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}
	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppID()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.PermissionDenied, "incorrect email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &ssov1.LoginResponse{Token: token}, nil
}

func (s Server) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed register user")
	}

	return &ssov1.RegisterResponse{UserID: userID}, nil
}

func (s Server) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if in.UserID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, in.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to check user is admin")
	}
	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
