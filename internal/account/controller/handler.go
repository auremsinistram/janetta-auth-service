package account

import (
	"auth-service/pkg/common"
	"auth-service/pkg/econst"
	"auth-service/pkg/tools"
	"net/http"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/response"
	"github.com/labstack/echo/v5"

	domain "auth-service/internal/account/domain"
)

func (c *controller) signUp(ctx *echo.Context) (err error) {
	var req domain.SignUpReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - signUp - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - signUp - #1")
	}

	if err = c.usecase.SignUp(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - signUp - #2")
	}

	return nil
}

func (c *controller) signUpConfirm(ctx *echo.Context) (err error) {
	var req domain.SignUpConfirmReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - signUpConfirm - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - signUpConfirm - #1")
	}

	if err = c.usecase.SignUpConfirm(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - signUpConfirm - #2")
	}

	return nil
}

func (c *controller) signIn(ctx *echo.Context) (err error) {
	var req domain.SignInReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - signIn - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - signIn - #1")
	}

	if err = c.usecase.SignIn(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - signIn - #2")
	}

	return nil
}

func (c *controller) signInConfirm(ctx *echo.Context) (err error) {
	var (
		req domain.SignInConfirmReq
		res domain.SignInConfirmRes
	)

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, res, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - signInConfirm - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - signInConfirm - #1")
	}

	if err = c.usecase.SignInConfirm(ctx.Request().Context(), &req, &res); err != nil {
		return errors.Wrap(err, "controller - signInConfirm - #2")
	}

	setCookie(ctx, res.AccessToken, &c.config.JWT.AccessToken.Cookie)
	setCookie(ctx, res.RefreshToken, &c.config.JWT.RefreshToken.Cookie)

	return nil
}

func (c *controller) refresh(ctx *echo.Context) (err error) {
	var (
		req domain.RefreshReq
		res domain.RefreshRes
	)

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, res, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - refresh - defer - #1")
		}
	}()

	// TODO: check cookie

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - refresh - #1")
	}

	if err = c.usecase.Refresh(ctx.Request().Context(), &req, &res); err != nil {
		return errors.Wrap(err, "controller - refresh - #2")
	}

	setCookie(ctx, res.AccessToken, &c.config.JWT.AccessToken.Cookie)
	setCookie(ctx, res.RefreshToken, &c.config.JWT.RefreshToken.Cookie)

	return nil
}

func (c *controller) logout(ctx *echo.Context) (err error) {
	var req domain.LogoutReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - logout - defer - #1")
		}
	}()

	// TODO: check cookie

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - logout - #1")
	}

	if err = c.usecase.Logout(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - logout - #2")
	}

	// TODO: clear cookie

	return nil
}

func (c *controller) passwordForgot(ctx *echo.Context) (err error) {
	var req domain.PasswordForgotReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - passwordForgot - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - passwordForgot - #1")
	}

	if err = c.usecase.PasswordForgot(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - passwordForgot - #2")
	}

	return nil
}

func (c *controller) passwordReset(ctx *echo.Context) (err error) {
	var req domain.PasswordResetReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - passwordReset - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - passwordReset - #1")
	}

	if err = c.usecase.PasswordReset(ctx.Request().Context(), &req); err != nil {
		return errors.Wrap(err, "controller - passwordReset - #2")
	}

	return nil
}

func (c *controller) passwordChange(ctx *echo.Context) (err error) {
	var req domain.PasswordChangeReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - passwordChange - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - passwordChange - #1")
	}

	sub, err := tools.GetSub(ctx)
	if err != nil {
		return errors.Wrap(err, "controller - passwordChange - #2")
	}

	if err = c.usecase.PasswordChange(ctx.Request().Context(), sub, &req); err != nil {
		return errors.Wrap(err, "controller - passwordChange - #3")
	}

	return nil
}

func (c *controller) passwordChangeConfirm(ctx *echo.Context) (err error) {
	var req domain.PasswordChangeConfirmReq

	defer func() {
		if e := response.SendOK(ctx, &common.Response{}, nil, err, econst.ErrRes, econst.ErrData); e != nil {
			err = errors.Wrap(e, "controller - passwordChangeConfirm - defer - #1")
		}
	}()

	if err = ctx.Bind(&req); err != nil {
		return errors.Wrap(err, "controller - passwordChangeConfirm - #1")
	}

	sub, err := tools.GetSub(ctx)
	if err != nil {
		return errors.Wrap(err, "controller - passwordChangeConfirm - #2")
	}

	if err = c.usecase.PasswordChangeConfirm(ctx.Request().Context(), sub, &req); err != nil {
		return errors.Wrap(err, "controller - passwordChangeConfirm - #3")
	}

	return nil
}

func setCookie(ctx *echo.Context, value string, cookie *domain.Cookie) {
	ctx.SetCookie(&http.Cookie{
		Name:     cookie.Name,
		Value:    value,
		Path:     cookie.Path,
		MaxAge:   int(cookie.MaxAge),
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		SameSite: http.SameSite(cookie.SameSite),
	})
}
