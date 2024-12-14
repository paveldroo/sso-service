package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	pb "github.com/paveldroo/sso-service/protos/sso"
	"github.com/paveldroo/sso-service/tests/suite"
	"github.com/stretchr/testify/require"
)

func TestGRPC(t *testing.T) {
	ts := suite.New(t)
	r := &pb.RegisterRequest{
		Email:    gofakeit.Email(),
		Password: suite.PassGen(),
	}
	resp, err := ts.Client.Register(context.Background(), r)
	require.NoError(t, err)

	fmt.Println(resp)
}
