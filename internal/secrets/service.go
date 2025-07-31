package secrets

import (
	"context"
	"log/slog"
	"github.com/ncostamagna/passit-back/adapters/cache"
)

type (
	
	Service interface {
		Create(ctx context.Context, secret *Secret) (string, error)
		Get(ctx context.Context, key string) (*Secret, error)
		Delete(ctx context.Context, key string) error
	}

	service struct {
		log  *slog.Logger
		cache cache.Cache
	}

)

func NewService(log *slog.Logger, cache cache.Cache) Service {
	return &service{log: log, cache: cache}
}

func (s *service) Create(ctx context.Context, secret *types.Secret) (string, error) {
	cache.Set(ctx, secret.Key, secret.Value, time.Duration(secret.Expiration)*time.Second)
}

func (s *service) Get(ctx context.Context, key string) (*types.Secret, error) {
	cache.Get(ctx, key)
}

func (s *service) Delete(ctx context.Context, key string) error {
	cache.Delete(ctx, key)
}