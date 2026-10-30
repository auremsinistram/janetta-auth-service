package account

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SignInAttempts int `mapstructure:"sign_in_attempts"`
	CodeLength     int `mapstructure:"code_length"`

	SignInAttemptsExp time.Duration `mapstructure:"sign_in_attempts_exp"`
	SignInBlockExp    time.Duration `mapstructure:"sign_in_block_exp"`
	SignUpExp         time.Duration `mapstructure:"sign_up_exp"`
	SignInExp         time.Duration `mapstructure:"sign_in_exp"`
	PasswordForgotExp time.Duration `mapstructure:"password_forgot_exp"`

	CodeCharset string `mapstructure:"code_charset"`

	JWT JWT `mapstructure:"jwt"`
}

type JWT struct {
	Issuer string `mapstructure:"issuer"`

	Audience     []string `mapstructure:"audience"`
	ValidMethods []string `mapstructure:"valid_methods"`

	AccessToken  Token `mapstructure:"access_token"`
	RefreshToken Token `mapstructure:"refresh_token"`
}

type Token struct {
	Exp time.Duration `mapstructure:"exp"`

	Algorithm string `mapstructure:"algorithm"`
	Secret    string `mapstructure:"secret"`

	Cookie Cookie `mapstructure:"cookie"`
}

type Cookie struct {
	Secure   bool `mapstructure:"secure"`
	HttpOnly bool `mapstructure:"http_only"`

	SameSite int `mapstructure:"same_site"`

	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`

	MaxAge time.Duration `mapstructure:"max_age"`
}

type Claims struct {
	jwt.RegisteredClaims

	TokenType string `json:"token_type"`
}

type AccessClaims struct {
	Claims

	// TODO: add access claims fields
}

type RefreshClaims struct {
	Claims
}

type SignUp struct {
	Code         string `json:"code"`
	PasswordHash string `json:"password_hash"`
}

type RefreshData struct {
	Sub string `redis:"sub"`

	CreatedAt  time.Time `redis:"created_at"`
	LastUsedAt time.Time `redis:"last_used_at"`
}

type Account struct {
	ID int

	UUID         string
	Phone        string
	PasswordHash string

	Metadata map[string]any

	CreatedAt time.Time

	UpdatedAt             *time.Time
	DeletedAt             *time.Time
	BlockedAt             *time.Time
	PasswordHashChangedAt *time.Time
}
