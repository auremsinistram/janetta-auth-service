package app

import (
	"auth-service/internal/account"
	"auth-service/internal/base"
	"auth-service/pkg/common"
	"auth-service/pkg/econst"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/logger"
	"github.com/auremsinistram/go-toolkit/middleware"
	"github.com/auremsinistram/go-toolkit/postgres"
	"github.com/auremsinistram/go-toolkit/redis"
	"github.com/auremsinistram/go-toolkit/response"
	"github.com/auremsinistram/go-toolkit/tools"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	config "auth-service/internal/account/domain"
	pkgtools "auth-service/pkg/tools"
)

func Run(configPath *string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}
	}()

	logger, err := logger.New()
	if err != nil {
		log.Panicln(errors.Wrap(err, "app - Run - #1"))
	}

	defer func() {
		logger.Sync()
		logger.Info("Service finished")
	}()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)

	viper.SetConfigFile(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		logger.Panic("app - Run - #2", zap.Error(err))
	}

	redis := redis.New()

	if err := redis.Connect(
		viper.GetString("redis.url"),
		viper.GetString("redis.username"),
		viper.GetString("redis.password"),
		viper.GetInt("redis.db"),
		viper.GetDuration("redis.timeout"),
	); err != nil {
		logger.Panic("app - Run - #3", zap.Error(err))
	}

	defer func() {
		if err := redis.Close(); err != nil {
			logger.Error("app - Run - defer - #1", zap.Error(err))
		}
	}()

	postgres := postgres.New()

	if err := postgres.Connect(
		viper.GetString("postgres.url"),
		viper.GetDuration("postgres.timeout"),
	); err != nil {
		logger.Panic("app - Run - #4", zap.Error(err))
	}

	defer func() {
		if err := postgres.Close(); err != nil {
			logger.Error("app - Run - defer - #2", zap.Error(err))
		}
	}()

	echo := echo.New()
	api := echo.Group("/api")

	shutting := &atomic.Bool{}

	errorsData, err := os.ReadFile("error.response.json")
	if err != nil {
		logger.Panic("app - Run - #5", zap.Error(err))
	}

	if err := base.Build(
		api,
		shutting,
		errorsData,
	); err != nil {
		logger.Panic("app - Run - #6", zap.Error(err))
	}

	var config config.Config

	if err := viper.UnmarshalKey("account", &config); err != nil {
		logger.Panic("app - Run - #7", zap.Error(err))
	}

	parser := jwt.NewParser(
		jwt.WithIssuer(config.JWT.Issuer),
		jwt.WithAudience(config.JWT.Audience...),
		jwt.WithValidMethods(config.JWT.ValidMethods),
	)

	authMiddleware := middleware.Auth(
		parser,
		config.JWT.AccessToken.Cookie.Name,
		middleware.StaticKey(config.JWT.AccessToken.Secret),
		claimsHandler(logger),
		errorHandler,
	)

	if err := account.Build(
		logger,
		&config,
		parser,
		api,
		redis,
		postgres,
		authMiddleware,
	); err != nil {
		logger.Panic("app - Run - #8", zap.Error(err))
	}

	echo.Use(
		middleware.RequestID(),
		middleware.RequestLogger(logger),
	)

	ongoingCtx, ongoingCancel := context.WithCancel(context.Background())
	defer ongoingCancel()

	server := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: echo,
		BaseContext: func(l net.Listener) context.Context {
			return ongoingCtx
		},
	}

	errorsChan := make(chan error, 1)
	defer close(errorsChan)

	go func() {
		logger.Info("Service started")

		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errorsChan <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		logger.Debug("Initiate closing")
	case err := <-errorsChan:
		logger.Panic("app - Run - #9", zap.Error(err))
	}

	if tools.GetenvBool("DEBUG", false) {
		cancel()
	}

	shutting.Store(true)

	logger.Debug("Wait readiness")

	<-time.After(viper.GetDuration("server.delays.readiness"))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), viper.GetDuration("server.delays.shutdown"))
	defer shutdownCancel()

	logger.Debug("Wait shutting")

	err = server.Shutdown(shutdownCtx)

	ongoingCancel()

	if err != nil {
		logger.Error("app - Run - #10", zap.Error(err))
		logger.Debug("Wait additional time")

		<-time.After(viper.GetDuration("server.delays.additional"))
	}
}

func claimsHandler(logger *zap.Logger) func(c *echo.Context, claims jwt.Claims) {
	return func(c *echo.Context, claims jwt.Claims) {
		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			return
		}

		if sub, ok := mapClaims["sub"].(string); ok {
			pkgtools.SetSub(c, sub)

			logger.Debug("app - claimsHandler - mapClaims", zap.String("sub", sub))
		}
	}
}

func errorHandler(c *echo.Context, err error) error {
	if errors.Is(err, middleware.ErrAuthHeaderMissing) {
		err = errors.Mark(err, econst.MissingAuthHeader)
	}

	if errors.Is(err, middleware.ErrAuthFormatInvalid) {
		err = errors.Mark(err, econst.InvalidAuthFormat)
	}

	// TODO: check others errors

	return response.SendError(c, &common.Response{}, err, econst.ErrRes, econst.ErrData)
}
