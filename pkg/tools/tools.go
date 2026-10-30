package tools

import (
	"github.com/auremsinistram/go-errors"
	"github.com/labstack/echo/v5"
)

const sub = "sub"

var (
	ErrSubEmpty   = errors.New("sub is empty")
	ErrSubInvalid = errors.New("sub is invalid")
)

func PasswordStrong(password string) bool {
	len := len(password)

	if len < 8 || len > 32 {
		return false
	}

	var (
		lower   bool
		upper   bool
		digit   bool
		special bool
	)

	for _, r := range password {
		if lower && upper && digit && special {
			return true
		}

		switch {
		case r >= 'a' && r <= 'z':
			lower = true
		case r >= 'A' && r <= 'Z':
			upper = true
		case r >= '0' && r <= '9':
			digit = true
		default:
			special = true
		}
	}

	return false
}

func GetSub(ctx *echo.Context) (string, error) {
	if s, ok := ctx.Get(sub).(string); ok {
		if s == "" {
			return "", ErrSubEmpty
		}

		return s, nil
	}

	return "", ErrSubInvalid
}

func SetSub(ctx *echo.Context, s string) {
	ctx.Set(sub, s)
}
