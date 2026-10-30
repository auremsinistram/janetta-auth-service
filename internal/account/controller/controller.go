package account

import (
	"github.com/labstack/echo/v5"

	domain "auth-service/internal/account/domain"
)

type controller struct {
	config *domain.Config

	api *echo.Group

	authMiddleware echo.MiddlewareFunc

	usecase domain.Usecase
}

func New(
	config *domain.Config,
	api *echo.Group,
	authMiddleware echo.MiddlewareFunc,
	usecase domain.Usecase,
) domain.Controller {
	return &controller{
		config:         config,
		api:            api,
		authMiddleware: authMiddleware,
		usecase:        usecase,
	}
}

func (c *controller) Expose() {
	// TODO: add description for swagger

	auth := c.api.Group("/auth")

	auth.POST("/signup", c.signUp)
	auth.POST("/signup/confirm", c.signUpConfirm)
	auth.POST("/signin", c.signIn)
	auth.POST("/signin/confirm", c.signInConfirm)
	auth.POST("/refresh", c.refresh)
	auth.POST("/logout", c.logout, c.authMiddleware)

	password := auth.Group("/password")

	password.POST("/forgot", c.passwordForgot)
	password.POST("/reset", c.passwordReset)
	password.POST("/change", c.passwordChange, c.authMiddleware)
	password.POST("/change/confirm", c.passwordChangeConfirm, c.authMiddleware)
}
