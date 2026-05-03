package secrets

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/entities"
)

type (
	Service interface {
		Create(ctx context.Context, secret *entities.Secret) (string, error)
		Get(ctx context.Context, key string) (*entities.Secret, error)
		Delete(ctx context.Context, key string) error
	}

	service struct {
		log   *slog.Logger
		cache cache.Cache
	}
)

func NewService(log *slog.Logger, cache cache.Cache) Service {
	return &service{log: log, cache: cache}
}

func (s *service) Create(ctx context.Context, secret *entities.Secret) (string, error) {
	s.log.Debug("[Create] secret - message:" + secret.Message + " | oneTime: " + strconv.FormatBool(secret.OneTime))
	key := uuid.New().String()
	value, err := secret.ToJSON()
	if err != nil {
		return "", err
	}

	if err := s.cache.Set(ctx, key, string(value), time.Duration(secret.Expiration)*time.Second); err != nil {
		s.log.Debug("[Create] Error secret: " + err.Error())
		return "", err
	}

	return key, nil
}

func (s *service) Get(ctx context.Context, key string) (*entities.Secret, error) {
	s.log.Debug("[Get] secret: " + key)
	secretJSON, err := s.cache.Get(ctx, key)
	if err != nil {
		s.log.Debug("[Error] secret not found")
		return nil, ErrSecretNotFound
	}

	var secret entities.Secret
	if err := secret.FromJSON([]byte(secretJSON)); err != nil {
		return nil, err
	}

	if secret.OneTime {
		if err := s.cache.Delete(ctx, key); err != nil {
			return nil, err
		}
	}

	s.log.Debug("[GET] secret: " + secret.Message)
	return &secret, nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, key)
}
