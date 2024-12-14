package suite

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/paveldroo/sso-service/internal/config"
	pb "github.com/paveldroo/sso-service/protos/sso"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestSuite struct {
	Client pb.AuthClient
}

func New(t *testing.T) TestSuite {
	t.Helper()

	cfg, err := config.ParseCfg("../../../config.yml")
	require.NoError(t, err)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := pb.NewAuthClient(conn)

	ts := TestSuite{Client: client}

	return ts
}

func PassGen() string {
	return gofakeit.Password(true, true, true, true, true, 5)
}
