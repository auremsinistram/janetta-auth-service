package account

import (
	"github.com/auremsinistram/go-toolkit/postgres"
	"github.com/auremsinistram/go-toolkit/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"

	controller "auth-service/internal/account/controller"
	domain "auth-service/internal/account/domain"
	cache "auth-service/internal/account/infrastructure/cache"
	limit "auth-service/internal/account/infrastructure/limit"
	repository "auth-service/internal/account/infrastructure/repository"
	usecase "auth-service/internal/account/usecase"
)

func Build(
	logger *zap.Logger,
	config *domain.Config,
	parser *jwt.Parser,
	api *echo.Group,
	redis *redis.Redis,
	postgres *postgres.Postgres,
	authMiddleware echo.MiddlewareFunc,
) error {
	controller.New(
		config,
		api,
		authMiddleware,
		usecase.New(
			logger,
			config,
			parser,
			cache.New(redis),
			limit.New(redis),
			repository.New(postgres),
		),
	).Expose()

	return nil
}
