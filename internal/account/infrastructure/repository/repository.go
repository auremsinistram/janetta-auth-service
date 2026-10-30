package account

import (
	"context"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/postgres"
	"github.com/jackc/pgx/v5"

	domain "auth-service/internal/account/domain"
)

type repository struct {
	postgres *postgres.Postgres
}

func New(postgres *postgres.Postgres) domain.Repository {
	return &repository{postgres: postgres}
}

func (c *repository) BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	tx, err := c.postgres.Pool.BeginTx(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "repository - BeginTx - #1")
	}

	return tx, nil
}

func (c *repository) Exists(ctx context.Context, tx pgx.Tx, phone string) (bool, error) {
	var exists bool

	// TODO: write query
	if err := tx.QueryRow(
		ctx,
		``,
		phone,
	).Scan(&exists); err != nil {
		return false, errors.Wrap(err, "repository - Exists - #1")
	}

	return exists, nil
}

func (c *repository) Account(ctx context.Context, tx pgx.Tx, phone string) (*domain.Account, error) {
	var account domain.Account

	// TODO: write query
	if err := tx.QueryRow(
		ctx,
		``,
		phone,
	).Scan(
		&account.ID,
		&account.UUID,
		&account.Phone,
		&account.PasswordHash,
		&account.Metadata,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.DeletedAt,
		&account.BlockedAt,
		&account.PasswordHashChangedAt,
	); err != nil {
		return nil, errors.Wrap(err, "repository - Account - #1")
	}

	return &account, nil
}

func (c *repository) Create(ctx context.Context, tx pgx.Tx, phone string, passwordHash string) error {
	// TODO: write exec and event
	if _, err := tx.Exec(
		ctx,
		``,
		phone,
		passwordHash,
	); err != nil {
		return errors.Wrap(err, "repository - Create - #1")
	}

	return nil
}

func (c *repository) UpdatePassword(ctx context.Context, tx pgx.Tx, phone string, passwordHash string) error {
	// TODO: write exec and event
	if _, err := tx.Exec(
		ctx,
		``,
		phone,
		passwordHash,
	); err != nil {
		return errors.Wrap(err, "repository - UpdatePassword - #1")
	}

	return nil
}
