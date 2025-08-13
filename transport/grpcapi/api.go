package grpcapi

import (
	"context"
	"log/slog"

	"github.com/ncostamagna/passit-back/internal/secrets"
	"github.com/ncostamagna/passit-back/internal/types"
	grpcPassit "github.com/ncostamagna/passit-proto/go/grpcPassit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type api struct {
	grpcPassit.UnimplementedPassitServer
	service secrets.Service
}

func New(app secrets.Service) *api {
	return &api{
		service: app,
	}
}

func (a *api) CreateSecret(ctx context.Context, in *grpcPassit.CreateSecretRequest) (*grpcPassit.CreateSecretResponse, error) {
	
	slog.Debug("Creating secret", "message", in.Message)
	
	secret := &types.Secret{
		OneTime: in.OneTime,
		Message: in.Message,
		Expiration: int32(in.Expiration),
	}

	key, err := a.service.Create(ctx, secret)
	if err != nil {
		slog.Error("Error creating secret", "error", err)
		return nil, status.Errorf(codes.Internal, "Error creating secret")
	}

	slog.Info("Secret created", "key", key)
	res := &grpcPassit.CreateSecretResponse{
		Key: key,
	}

	return res, nil
}

func (a *api) GetSecret(ctx context.Context, in *grpcPassit.GetSecretRequest) (*grpcPassit.GetSecretResponse, error) {	
	
	slog.Debug("Getting secret", "key", in.Key)

	secret, err := a.service.Get(ctx, in.Key)
	if err != nil {
		if err == secrets.ErrSecretNotFound {
			return nil, status.Errorf(codes.NotFound, "Secret not found")
		}

		slog.Error("Error getting secret", "error", err)
		return nil, status.Errorf(codes.Internal, "Error getting secret")
	}

	res := &grpcPassit.GetSecretResponse{
		Message: secret.Message,
	}

	return res, nil
}


func (a *api) Register(s grpc.ServiceRegistrar) {
	grpcPassit.RegisterPassitServer(s, a)
}