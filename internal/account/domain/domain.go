package account

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Controller interface {
	Expose()
}

type Cache interface {
	GetSignUp(ctx context.Context, key string) (*SignUp, error)
	SetSignUp(ctx context.Context, key string, signUp *SignUp, exp time.Duration) error
	GetCode(ctx context.Context, key string) (string, error)
	SetCode(ctx context.Context, key string, code string, exp time.Duration) error
	AddRefreshData(ctx context.Context, jti string, refreshData *RefreshData, now time.Time, exp time.Time) error
	Delete(ctx context.Context, key string) error
}

type Limit interface {
	Blocked(ctx context.Context, key string) (bool, error)
	Attempts(ctx context.Context, key string) (int, error)
	Block(ctx context.Context, key string, exp time.Duration) error
	MakeAttempt(ctx context.Context, key string, exp time.Duration) error
}

type Repository interface {
	BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error)
	Exists(ctx context.Context, tx pgx.Tx, phone string) (bool, error)
	Account(ctx context.Context, tx pgx.Tx, phone string) (*Account, error)
	Create(ctx context.Context, tx pgx.Tx, phone string, passwordHash string) error
	UpdatePassword(ctx context.Context, tx pgx.Tx, phone string, passwordHash string) error
}

type Usecase interface {
	SignUp(ctx context.Context, req *SignUpReq) error
	SignUpConfirm(ctx context.Context, req *SignUpConfirmReq) error
	SignIn(ctx context.Context, req *SignInReq) error
	SignInConfirm(ctx context.Context, req *SignInConfirmReq, res *SignInConfirmRes) error
	Logout(ctx context.Context, req *LogoutReq) error
	Refresh(ctx context.Context, req *RefreshReq, res *RefreshRes) error
	PasswordForgot(ctx context.Context, req *PasswordForgotReq) error
	PasswordReset(ctx context.Context, req *PasswordResetReq) error
	PasswordChange(ctx context.Context, uuid string, req *PasswordChangeReq) error
	PasswordChangeConfirm(ctx context.Context, uuid string, req *PasswordChangeConfirmReq) error
}
