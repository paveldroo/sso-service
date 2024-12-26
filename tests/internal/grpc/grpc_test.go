package grpc_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	ssov1 "github.com/paveldroo/sso-service/protos/gen/go/sso"
	"github.com/paveldroo/sso-service/tests/suite"
	"github.com/stretchr/testify/require"
)

func TestGRPC(t *testing.T) {
	ts := suite.New(t)
	email := gofakeit.Email()
	password := suite.PassGen()
	appID := gofakeit.Uint32()

	r := &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	}
	resp, err := ts.Client.Register(context.Background(), r)
	require.NoError(t, err)
	require.NotEmpty(t, resp.UserID)

	rl := &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppID:    appID,
	}
	loginResp, err := ts.Client.Login(context.Background(), rl)
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Token)

	ar := &ssov1.IsAdminRequest{
		UserID: resp.UserID,
	}
	adminResp, err := ts.Client.IsAdmin(context.Background(), ar)
	require.NoError(t, err)
	require.Equal(t, false, adminResp.IsAdmin)
}
