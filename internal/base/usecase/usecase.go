package base

import (
	"context"
	"sync/atomic"

	domain "auth-service/internal/base/domain"
)

type usecase struct {
	shutting *atomic.Bool

	errorsRes *domain.ErrorsRes
}

func New(
	shutting *atomic.Bool,
	errorsRes *domain.ErrorsRes,
) domain.Usecase {
	return &usecase{
		shutting:  shutting,
		errorsRes: errorsRes,
	}
}

func (u *usecase) Shutting(ctx context.Context) (bool, error) {
	return u.shutting.Load(), nil
}

func (u *usecase) Errors(ctx context.Context, res *domain.ErrorsRes) error {
	*res = *u.errorsRes

	return nil
}
