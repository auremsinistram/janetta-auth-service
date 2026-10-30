package base

import (
	"encoding/json"
	"sync/atomic"

	"github.com/auremsinistram/go-errors"
	"github.com/labstack/echo/v5"

	controller "auth-service/internal/base/controller"
	domain "auth-service/internal/base/domain"
	usecase "auth-service/internal/base/usecase"
)

func Build(
	api *echo.Group,
	shutting *atomic.Bool,
	errorsData []byte,
) error {
	var errorsRes domain.ErrorsRes

	if err := json.Unmarshal(errorsData, &errorsRes); err != nil {
		return errors.Wrap(err, "builder - Build - #1")
	}

	controller.New(
		api,
		usecase.New(
			shutting,
			&errorsRes,
		),
	).Expose()

	return nil
}
