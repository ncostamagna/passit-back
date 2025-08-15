package instance

import (
	"log/slog"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/secrets"
)

func NewSecretService(cache cache.Cache, logger *slog.Logger) secrets.Service {
	return secrets.NewService(logger, cache)
}
