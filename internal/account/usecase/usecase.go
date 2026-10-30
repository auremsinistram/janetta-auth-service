package account

import (
	"auth-service/pkg/econst"
	"context"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/middleware"
	"github.com/auremsinistram/go-toolkit/tools"
	"github.com/auremsinistram/go-toolkit/transaction"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	domain "auth-service/internal/account/domain"
)

const (
	accessTokenType     = "access"
	refreshTokenType    = "refresh"
	signUpPhone         = "sign-up:phone:"
	signInPhone         = "sign-in:phone:"
	passwordForgotPhone = "password-forgot:phone:"
)

type signFunc func(token *jwt.Token) (string, error)

type usecase struct {
	logger *zap.Logger

	config *domain.Config

	parser *jwt.Parser

	cache domain.Cache

	limit domain.Limit

	repository domain.Repository
}

func New(
	logger *zap.Logger,
	config *domain.Config,
	parser *jwt.Parser,
	cache domain.Cache,
	limit domain.Limit,
	repository domain.Repository,
) domain.Usecase {
	return &usecase{
		logger:     logger,
		config:     config,
		parser:     parser,
		cache:      cache,
		limit:      limit,
		repository: repository,
	}
}

func (u *usecase) SignUp(ctx context.Context, req *domain.SignUpReq) (err error) {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - SignUp - #1")
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUp - #2")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - SignUp - defer - #1")
		}
	}()

	exists, err := u.repository.Exists(ctx, tx, req.Phone)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUp - #3")
	}

	if exists {
		return errors.NewWithCode(econst.PhoneAlreadyExists, "exists")
	}

	key := signUpPhone + req.Phone

	cachedSignUp, err := u.cache.GetSignUp(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUp - #4")
	}

	if cachedSignUp != nil {
		return errors.NewWithCode(econst.CodeAlreadySent, "cached sign up is not nil")
	}

	code, err := tools.GenerateCode(u.config.CodeCharset, u.config.CodeLength)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUp - #5")
	}

	u.logger.Debug("usecase - SignUp - GenerateCode", zap.String("code", code))

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUp - #6")
	}

	signUp := &domain.SignUp{
		Code:         code,
		PasswordHash: string(bytes),
	}

	if err := u.cache.SetSignUp(ctx, key, signUp, u.config.SignUpExp); err != nil {
		return errors.Wrap(err, "usecase - SignUp - #7")
	}

	// TODO: send sms

	return nil
}

func (u *usecase) SignUpConfirm(ctx context.Context, req *domain.SignUpConfirmReq) error {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - SignUpConfirm - #1")
	}

	key := signUpPhone + req.Phone

	signUp, err := u.cache.GetSignUp(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUpConfirm - #2")
	}

	if signUp == nil {
		return errors.NewWithCode(econst.CodeNotFound, "sign up is nil")
	}

	if req.Code != signUp.Code {
		return errors.NewWithCode(econst.InvalidCode, "invalid code")
	}

	if err := u.cache.Delete(ctx, key); err != nil {
		u.logger.Error("usecase - SignUpConfirm - #3", zap.Error(err))
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - SignUpConfirm - #4")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - SignUpConfirm - defer - #1")
		}
	}()

	if err := u.repository.Create(ctx, tx, req.Phone, signUp.PasswordHash); err != nil {
		return errors.Wrap(err, "usecase - SignUpConfirm - #5")
	}

	return nil
}

func (u *usecase) SignIn(ctx context.Context, req *domain.SignInReq) (err error) {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - SignIn - #1")
	}

	key := signInPhone + req.Phone

	blocked, err := u.limit.Blocked(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - SignIn - #2")
	}

	if blocked {
		return errors.NewWithCode(econst.TooManyAttempts, "blocked")
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - SignIn - #3")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - SignIn - defer - #1")
		}
	}()

	account, err := u.repository.Account(ctx, tx, req.Phone)
	if err != nil {
		return errors.Wrap(err, "usecase - SignIn - #4")
	}

	if account.BlockedAt != nil {
		return errors.NewWithCode(econst.AccountBlocked, "blocked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		attempts, err := u.limit.Attempts(ctx, key)
		if err != nil {
			return errors.Wrap(err, "usecase - SignIn - #5")
		}

		if attempts < u.config.SignInAttempts {
			if err := u.limit.MakeAttempt(ctx, key, u.config.SignInAttemptsExp); err != nil {
				return errors.Wrap(err, "usecase - SignIn - #6")
			}

			return errors.NewWithCode(econst.InvalidCredentials, "invalid credentials")
		}

		if err := u.limit.Block(ctx, key, u.config.SignUpExp); err != nil {
			return errors.Wrap(err, "usecase - SignIn - #7")
		}

		return errors.NewWithCode(econst.TooManyAttempts, "blocked")
	}

	code, err := tools.GenerateCode(u.config.CodeCharset, u.config.CodeLength)
	if err != nil {
		return errors.Wrap(err, "usecase - SignIn - #8")
	}

	u.logger.Debug("usecase - SignIn - GenerateCode", zap.String("code", code))

	if err := u.cache.SetCode(ctx, key, code, u.config.SignInExp); err != nil {
		return errors.Wrap(err, "usecase - SignIn - #9")
	}

	// TODO: send sms

	return nil
}

func (u *usecase) SignInConfirm(ctx context.Context, req *domain.SignInConfirmReq, res *domain.SignInConfirmRes) error {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #1")
	}

	key := signInPhone + req.Phone

	code, err := u.cache.GetCode(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #2")
	}

	if code == "" {
		return errors.NewWithCode(econst.CodeNotFound, "code is empty")
	}

	if req.Code != code {
		return errors.NewWithCode(econst.InvalidCode, "invalid code")
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #3")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - SignInConfirm - defer - #1")
		}
	}()

	account, err := u.repository.Account(ctx, tx, req.Phone)
	if err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #4")
	}

	if account.BlockedAt != nil {
		return errors.NewWithCode(econst.AccountBlocked, "blocked")
	}

	now := time.Now()

	accessToken, err := generateToken(
		createAccessClaims(
			uuid.NewString(),
			account.UUID,
			u.config.JWT.Issuer,
			u.config.JWT.Audience,
			u.config.JWT.AccessToken.Exp,
			now,
		),
		u.config.JWT.AccessToken.Algorithm,
		sign(u.config.JWT.AccessToken.Secret),
	)
	if err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #5")
	}

	refreshJTI := uuid.NewString()

	refreshToken, err := generateToken(
		createRefreshClaims(
			refreshJTI,
			account.UUID,
			u.config.JWT.Issuer,
			u.config.JWT.Audience,
			u.config.JWT.RefreshToken.Exp,
			now,
		),
		u.config.JWT.RefreshToken.Algorithm,
		sign(u.config.JWT.RefreshToken.Secret),
	)
	if err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #6")
	}

	refreshData := &domain.RefreshData{
		Sub:        account.UUID,
		CreatedAt:  now,
		LastUsedAt: now,
	}

	if err := u.cache.AddRefreshData(
		ctx,
		refreshJTI,
		refreshData,
		now,
		now.Add(u.config.JWT.RefreshToken.Exp),
	); err != nil {
		return errors.Wrap(err, "usecase - SignInConfirm - #7")
	}

	if err := u.cache.Delete(ctx, key); err != nil {
		u.logger.Error("usecase - SignInConfirm - #8", zap.Error(err))
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	return nil
}

func (u *usecase) Refresh(ctx context.Context, req *domain.RefreshReq, res *domain.RefreshRes) error {
	var claims domain.RefreshClaims

	if _, err := u.parser.ParseWithClaims(
		req.RefreshToken,
		&claims,
		middleware.StaticKey(u.config.JWT.RefreshToken.Secret),
	); err != nil {
		// TODO: check others errors

		return errors.Wrap(err, "usecase - Refresh - #1")
	}

	if claims.TokenType != refreshTokenType {
		return errors.NewWithCode(econst.Unauthorized, "token type is invalid")
	}

	// TODO: implement

	return nil
}

func (u *usecase) Logout(ctx context.Context, req *domain.LogoutReq) error {
	// TODO: implement

	return errors.New("unimplemented method")
}

func (u *usecase) PasswordForgot(ctx context.Context, req *domain.PasswordForgotReq) error {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #1")
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #2")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - PasswordForgot - defer - #1")
		}
	}()

	exists, err := u.repository.Exists(ctx, tx, req.Phone)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #3")
	}

	if !exists {
		return errors.NewWithCode(econst.PhoneNotFound, "not exists")
	}

	key := passwordForgotPhone + req.Phone

	cachedCode, err := u.cache.GetCode(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #4")
	}

	if cachedCode != "" {
		return errors.NewWithCode(econst.CodeAlreadySent, "cached code is not empty")
	}

	code, err := tools.GenerateCode(u.config.CodeCharset, u.config.CodeLength)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #5")
	}

	u.logger.Debug("usecase - PasswordForgot - GenerateCode", zap.String("code", code))

	if err := u.cache.SetCode(ctx, key, code, u.config.PasswordForgotExp); err != nil {
		return errors.Wrap(err, "usecase - PasswordForgot - #6")
	}

	// TODO: send sms

	return nil
}

func (u *usecase) PasswordReset(ctx context.Context, req *domain.PasswordResetReq) error {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - PasswordReset - #1")
	}

	key := passwordForgotPhone + req.Phone

	code, err := u.cache.GetCode(ctx, key)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordReset - #2")
	}

	if code == "" {
		return errors.NewWithCode(econst.CodeNotFound, "code is empty")
	}

	if req.Code != code {
		return errors.NewWithCode(econst.InvalidCode, "invalid code")
	}

	options := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	tx, err := u.repository.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(err, "usecase - PasswordReset - #3")
	}

	defer func() {
		if e := transaction.HandleTx(ctx, tx, err); e != nil {
			err = errors.Wrap(e, "usecase - PasswordReset - defer - #1")
		}
	}()

	if err := u.repository.UpdatePassword(ctx, tx, req.Phone, req.Password); err != nil {
		return errors.Wrap(err, "usecase - PasswordReset - #4")
	}

	if err := u.cache.Delete(ctx, key); err != nil {
		u.logger.Error("usecase - PasswordReset - #5", zap.Error(err))
	}

	return nil
}

func (u *usecase) PasswordChange(ctx context.Context, uuid string, req *domain.PasswordChangeReq) error {
	if err := req.Validate(); err != nil {
		return errors.Wrap(err, "usecase - PasswordChange - #1")
	}

	// TODO: implement

	return errors.New("unimplemented method")
}

func (u *usecase) PasswordChangeConfirm(ctx context.Context, uuid string, req *domain.PasswordChangeConfirmReq) error {
	// TODO: implement

	return errors.New("unimplemented method")
}

func createAccessClaims(
	id string,
	subject string,
	issuer string,
	audience []string,
	exp time.Duration,
	time time.Time,
) domain.AccessClaims {
	return domain.AccessClaims{
		Claims: domain.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				Subject:   subject,
				Issuer:    issuer,
				Audience:  audience,
				IssuedAt:  jwt.NewNumericDate(time),
				ExpiresAt: jwt.NewNumericDate(time.Add(exp)),
			},
			TokenType: accessTokenType,
		},
	}
}

func createRefreshClaims(
	id string,
	subject string,
	issuer string,
	audience []string,
	exp time.Duration,
	time time.Time,
) domain.RefreshClaims {
	return domain.RefreshClaims{
		Claims: domain.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				Subject:   subject,
				Issuer:    issuer,
				Audience:  audience,
				IssuedAt:  jwt.NewNumericDate(time),
				ExpiresAt: jwt.NewNumericDate(time.Add(exp)),
			},
			TokenType: refreshTokenType,
		},
	}
}

func sign(secret string) signFunc {
	return func(token *jwt.Token) (string, error) {
		return token.SignedString([]byte(secret))
	}
}

func generateToken(
	claims jwt.Claims,
	algorithm string,
	sign signFunc,
	opts ...jwt.TokenOption,
) (string, error) {
	signingMethod := jwt.GetSigningMethod(algorithm)

	if signingMethod == nil {
		return "", errors.New("signing method is nil")
	}

	token := jwt.NewWithClaims(signingMethod, claims, opts...)

	signedString, err := sign(token)
	if err != nil {
		return "", errors.Wrap(err, "usecase - generateToken - #1")
	}

	return signedString, nil
}
