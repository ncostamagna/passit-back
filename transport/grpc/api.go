package grpcapi

import (
	"github.com/ncostamagna/passit-back/internal/secrets"
	grpcPassit "github.com/ncostamagna/passit-proto/go/grpcPassit"
	"google.golang.org/grpc"
)

type api struct {
	service secrets.Service
}

func New(app secrets.Service) *api {
	return &api{
		service: app,
	}
}

func (a *api) Register(s *grpc.ServiceRegistrar) {
	grpcPassit.RegisterPassitServiceServer(s, a)
}