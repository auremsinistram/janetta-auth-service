package base

import (
	"auth-service/pkg/common"
	"auth-service/pkg/econst"
	"net/http"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/response"
	"github.com/labstack/echo/v5"

	domain "auth-service/internal/base/domain"
)

func (c *controller) health(ctx *echo.Context) error {
	shutting, err := c.usecase.Shutting(ctx.Request().Context())
	if err != nil {
		return errors.Wrap(err, "controller - health - #1")
	}

	code := http.StatusOK

	if shutting {
		code = http.StatusServiceUnavailable
	}

	if err := ctx.NoContent(code); err != nil {
		return errors.Wrap(err, "controller - health - #2")
	}

	return nil
}

func (c *controller) errors(ctx *echo.Context) (err error) {
	var res domain.ErrorsRes

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, res, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - errors - defer - #1")
		}
	}()

	if err = c.usecase.Errors(ctx.Request().Context(), &res); err != nil {
		return errors.Wrap(err, "controller - errors - #1")
	}

	return nil
}
