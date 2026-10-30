package base

import "context"

type Controller interface {
	Expose()
}

type Usecase interface {
	Shutting(ctx context.Context) (bool, error)
	Errors(ctx context.Context, res *ErrorsRes) error
}
