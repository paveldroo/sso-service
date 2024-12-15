package grpc_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	pb "github.com/paveldroo/sso-service/protos/sso"
	"github.com/paveldroo/sso-service/tests/suite"
	"github.com/stretchr/testify/require"
)

func TestGRPC(t *testing.T) {
	ts := suite.New(t)
	email := gofakeit.Email()
	password := suite.PassGen()

	r := &pb.RegisterRequest{
		Email:    email,
		Password: password,
	}
	resp, err := ts.Client.Register(context.Background(), r)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Message)

	rl := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}
	loginResp, err := ts.Client.Login(context.Background(), rl)
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Token)

	ar := &pb.IsAdminRequest{
		Email: email,
	}
	adminResp, err := ts.Client.IsAdmin(context.Background(), ar)
	require.NoError(t, err)
	require.Equal(t, false, adminResp.IsAdmin)
}
