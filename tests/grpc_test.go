package tests

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/secrets"
	"github.com/ncostamagna/passit-back/transport/grpcapi"
	grpcPassit "github.com/ncostamagna/passit-proto/go/grpcPassit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func redisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "localhost:6381"
}

func redisPass() string {
	if v := os.Getenv("REDIS_PASS"); v != "" {
		return v
	}
	return "admin"
}

func newTestServer(t *testing.T) (*bufconn.Listener, func()) {
	t.Helper()

	c := cache.NewCache(redisAddr(), redisPass(), 0)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	svc := secrets.NewService(logger, c)
	api := grpcapi.New(svc)

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	api.Register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("bufconn server stopped: %v", err)
		}
	}()

	return lis, func() {
		srv.GracefulStop()
		lis.Close()
	}
}

func newTestClient(t *testing.T, lis *bufconn.Listener) grpcPassit.PassitClient {
	t.Helper()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial bufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	return grpcPassit.NewPassitClient(conn)
}

func TestCreateSecret(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	resp, err := client.CreateSecret(context.Background(), &grpcPassit.CreateSecretRequest{
		Message:    "hello world",
		OneTime:    false,
		Expiration: 60,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.GetKey())
}

func TestGetSecret(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	ctx := context.Background()

	createResp, err := client.CreateSecret(ctx, &grpcPassit.CreateSecretRequest{
		Message:    "super secret",
		OneTime:    false,
		Expiration: 60,
	})
	require.NoError(t, err)

	getResp, err := client.GetSecret(ctx, &grpcPassit.GetSecretRequest{Key: createResp.GetKey()})
	require.NoError(t, err)
	assert.Equal(t, "super secret", getResp.GetMessage())
}

func TestGetSecret_NotFound(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	_, err := client.GetSecret(context.Background(), &grpcPassit.GetSecretRequest{Key: "non-existent-key"})
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGetSecret_OneTime(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	ctx := context.Background()

	createResp, err := client.CreateSecret(ctx, &grpcPassit.CreateSecretRequest{
		Message:    "burn after reading",
		OneTime:    true,
		Expiration: 60,
	})
	require.NoError(t, err)
	key := createResp.GetKey()

	getResp, err := client.GetSecret(ctx, &grpcPassit.GetSecretRequest{Key: key})
	require.NoError(t, err)
	assert.Equal(t, "burn after reading", getResp.GetMessage())

	_, err = client.GetSecret(ctx, &grpcPassit.GetSecretRequest{Key: key})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGetSecret_Expired(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	ctx := context.Background()

	createResp, err := client.CreateSecret(ctx, &grpcPassit.CreateSecretRequest{
		Message:    "short lived",
		OneTime:    false,
		Expiration: 1,
	})
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = client.GetSecret(ctx, &grpcPassit.GetSecretRequest{Key: createResp.GetKey()})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGetSecret_MultipleReads_NotOneTime(t *testing.T) {
	lis, stop := newTestServer(t)
	defer stop()
	client := newTestClient(t, lis)

	ctx := context.Background()

	createResp, err := client.CreateSecret(ctx, &grpcPassit.CreateSecretRequest{
		Message:    "persistent secret",
		OneTime:    false,
		Expiration: 60,
	})
	require.NoError(t, err)
	key := createResp.GetKey()

	for i := range 3 {
		resp, err := client.GetSecret(ctx, &grpcPassit.GetSecretRequest{Key: key})
		require.NoErrorf(t, err, "read %d", i+1)
		assert.Equalf(t, "persistent secret", resp.GetMessage(), "read %d", i+1)
	}
}
