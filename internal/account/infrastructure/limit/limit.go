package account

import (
	"context"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/redis"

	goredis "github.com/redis/go-redis/v9"

	domain "auth-service/internal/account/domain"
)

const (
	blocked  = "blocked:"
	attempts = "attempts:"
)

type limit struct {
	redis *redis.Redis
}

func New(redis *redis.Redis) domain.Limit {
	return &limit{redis: redis}
}

func (c *limit) Blocked(ctx context.Context, key string) (bool, error) {
	count, err := c.redis.Client.Exists(ctx, blocked+key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return false, nil
		}

		return false, errors.Wrap(err, "limit - Blocked - #1")
	}

	return count > 0, nil
}

func (c *limit) Attempts(ctx context.Context, key string) (int, error) {
	attempts, err := c.redis.Client.Get(ctx, attempts+key).Int()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, nil
		}

		return 0, errors.Wrap(err, "limit - Attempts - #1")
	}

	return attempts, nil
}

func (c *limit) Block(ctx context.Context, key string, exp time.Duration) error {
	if err := c.redis.Client.Set(ctx, blocked+key, true, exp).Err(); err != nil {
		return errors.Wrap(err, "limit - Block - #1")
	}

	return nil
}

func (c *limit) MakeAttempt(ctx context.Context, phone string, exp time.Duration) error {
	key := attempts + phone

	pipe := c.redis.Client.TxPipeline()

	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, exp)

	if cmds, err := pipe.Exec(ctx); err != nil {
		for i, cmd := range cmds {
			if e := cmd.Err(); e != nil {
				return errors.Wrapf(e, "limit - MakeAttempt - cmd - #%d", i+1)
			}
		}

		return errors.Wrap(err, "limit - MakeAttempt - #1")
	}

	return nil
}
