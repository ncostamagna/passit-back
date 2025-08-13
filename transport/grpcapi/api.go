package grpcapi

import (
	"context"
	"log/slog"

	"github.com/ncostamagna/passit-back/internal/secrets"
	//"github.com/ncostamagna/passit-back/internal/types"
	grpcPassit "github.com/ncostamagna/passit-proto/go/grpcPassit"
	"google.golang.org/grpc"
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
	
	slog.Info("Creating secret", "message", in.Message)
	
	res := &grpcPassit.CreateSecretResponse{
		ClientSecret: "123",
	}

	return res, nil
}


func (a *api) Register(s grpc.ServiceRegistrar) {
	grpcPassit.RegisterPassitServer(s, a)
}