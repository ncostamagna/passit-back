package grpcapi

import (
	"context"
	"log/slog"

	"github.com/ncostamagna/passit-back/internal/entities"
	"github.com/ncostamagna/passit-back/internal/secrets"
	grpcPassit "github.com/ncostamagna/passit-proto/go/grpcPassit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API struct {
	grpcPassit.UnimplementedPassitServer
	service secrets.Service
}

func New(app secrets.Service) *API {
	return &API{
		service: app,
	}
}

func (a *API) CreateSecret(ctx context.Context, in *grpcPassit.CreateSecretRequest) (*grpcPassit.CreateSecretResponse, error) {

	slog.Debug("Creating secret", "message", in.GetMessage())

	secret := &entities.Secret{
		OneTime:    in.GetOneTime(),
		Message:    in.GetMessage(),
		Expiration: in.GetExpiration(),
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

func (a *API) GetSecret(ctx context.Context, in *grpcPassit.GetSecretRequest) (*grpcPassit.GetSecretResponse, error) {

	slog.Debug("Getting secret", "key", in.GetKey())

	secret, err := a.service.Get(ctx, in.GetKey())
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

func (a *API) Register(s grpc.ServiceRegistrar) {
	grpcPassit.RegisterPassitServer(s, a)
}
