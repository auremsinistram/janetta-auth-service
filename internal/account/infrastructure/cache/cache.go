package account

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/redis"

	goredis "github.com/redis/go-redis/v9"

	domain "auth-service/internal/account/domain"
)

const (
	refresh  = "refresh:"
	sessions = "sessions:"
)

type cache struct {
	redis *redis.Redis
}

func New(redis *redis.Redis) domain.Cache {
	return &cache{redis: redis}
}

func (c *cache) GetSignUp(ctx context.Context, key string) (*domain.SignUp, error) {
	data, err := c.redis.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "cache - GetSignUp - #1")
	}

	var signUp domain.SignUp

	if err := json.Unmarshal(data, &signUp); err != nil {
		return nil, errors.Wrap(err, "cache - GetSignUp - #2")
	}

	return &signUp, nil
}

func (c *cache) SetSignUp(
	ctx context.Context,
	key string,
	signUp *domain.SignUp,
	exp time.Duration,
) error {
	if signUp == nil {
		return errors.New("sign up is nil")
	}

	data, err := json.Marshal(signUp)
	if err != nil {
		return errors.Wrap(err, "cache - SetSignUp - #1")
	}

	if err := c.redis.Client.Set(ctx, key, data, exp).Err(); err != nil {
		return errors.Wrap(err, "cache - SetSignUp - #2")
	}

	return nil
}

func (c *cache) GetCode(ctx context.Context, key string) (string, error) {
	code, err := c.redis.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", nil
		}

		return "", errors.Wrap(err, "cache - GetCode - #1")
	}

	return code, nil
}

func (c *cache) SetCode(ctx context.Context, key string, code string, exp time.Duration) error {
	if err := c.redis.Client.Set(ctx, key, code, exp).Err(); err != nil {
		return errors.Wrap(err, "cache - SetCode - #1")
	}

	return nil
}

func (c *cache) AddRefreshData(
	ctx context.Context,
	jti string,
	refreshData *domain.RefreshData,
	now time.Time,
	exp time.Time,
) error {
	if refreshData == nil {
		return errors.New("refresh data is nil")
	}

	sub := strings.TrimSpace(refreshData.Sub)

	if sub == "" {
		return errors.New("refresh data's sub is empty")
	}

	refreshKey := refresh + jti
	sessionsKey := sessions + sub

	cutoff := strconv.FormatInt(now.Unix(), 10)

	member := goredis.Z{
		Score:  float64(exp.Unix()),
		Member: jti,
	}

	pipe := c.redis.Client.TxPipeline()

	pipe.HSet(ctx, refreshKey, refreshData)
	pipe.ExpireAt(ctx, refreshKey, exp)
	pipe.ZRemRangeByScore(ctx, sessionsKey, "-inf", cutoff)
	pipe.ZAdd(ctx, sessionsKey, member)
	pipe.ExpireAt(ctx, sessionsKey, exp)

	if cmds, err := pipe.Exec(ctx); err != nil {
		for i, cmd := range cmds {
			if e := cmd.Err(); e != nil {
				return errors.Wrapf(e, "cache - AddRefreshData - cmd - #%d", i+1)
			}
		}

		return errors.Wrap(err, "cache - AddRefreshData - #1")
	}

	return nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	if err := c.redis.Client.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "cache - Delete - #1")
	}

	return nil
}
