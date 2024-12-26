package suite

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/paveldroo/sso-service/internal/config"
	ssov1 "github.com/paveldroo/sso-service/protos/gen/go/sso"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestSuite struct {
	Client ssov1.AuthClient
}

func New(t *testing.T) TestSuite {
	t.Helper()

	cfg := config.MustLoadPath("../../../config/config.yml")
	addr := fmt.Sprintf("localhost:%d", cfg.GRPC.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := ssov1.NewAuthClient(conn)

	ts := TestSuite{Client: client}

	return ts
}

func PassGen() string {
	return gofakeit.Password(true, true, true, true, true, 5)
}
