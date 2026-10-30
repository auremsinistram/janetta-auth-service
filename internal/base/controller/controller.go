package base

import (
	"github.com/labstack/echo/v5"

	domain "auth-service/internal/base/domain"
)

type controller struct {
	api *echo.Group

	usecase domain.Usecase
}

func New(
	api *echo.Group,
	usecase domain.Usecase,
) domain.Controller {
	return &controller{
		api:     api,
		usecase: usecase,
	}
}

func (c *controller) Expose() {
	// TODO: add description for swagger

	c.api.GET("/health", c.health)
	c.api.GET("/errors", c.errors)
}
