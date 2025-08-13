package secrets

import (
	"context"
	"log/slog"
	"time"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/types"
	"github.com/google/uuid"
)

type (
	
	Service interface {
		Create(ctx context.Context, secret *types.Secret) (string, error)
		Get(ctx context.Context, key string) (*types.Secret, error)
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
	key := uuid.New().String()
	value, err := secret.ToJSON()
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, key, string(value), time.Duration(secret.Expiration)*time.Second); err != nil {
		return "", err
	}

	return key, nil
}

func (s *service) Get(ctx context.Context, key string) (*types.Secret, error) {
	secretJson, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var secret types.Secret
	if err := secret.FromJSON([]byte(secretJson)); err != nil {
		return nil, err
	}

	if secret.OneTime {
		if err := s.cache.Delete(ctx, key); err != nil {
			return nil, err
		}
	}

	return &secret, nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, key)
}