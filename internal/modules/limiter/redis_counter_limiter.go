package limiter

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

func newRedisCounterLimiter(interval, maxCount int) *redisCounterLimiter {
	counterLuaScript := ` 
    local count = redis.call('incr',KEYS[1]);
	if tonumber(count) == 1 then
		redis.call('expire', KEYS[1], ARGV[1])
	end
	return count`

	return &redisCounterLimiter{
		counterLuaScript: counterLuaScript,

		interval: interval,
		maxCount: maxCount,
	}
}

type redisCounterLimiter struct {
	counterLuaScript string

	interval int
	maxCount int

	evalSha string
}

func (l *redisCounterLimiter) initScript(ctx context.Context, rdb *redis.Client) error {
	evalSha, err := rdb.ScriptLoad(ctx, l.counterLuaScript).Result()
	if err != nil {
		return err
	}

	l.evalSha = evalSha
	return nil
}

func (l *redisCounterLimiter) IsActionAllow(ctx context.Context, rdb *redis.Client, ip string) (bool, int64, error) {
	res, err2 := rdb.EvalSha(ctx, l.evalSha, []string{ip}, l.interval).Result()
	if err2 != nil {
		return false, 0, err2
	}

	var count int64
	if v, ok := res.(int64); !ok {
		return false, 0, errors.New("redis result count Type Assertion int64 fail")
	} else {
		count = v
	}

	if count > int64(l.maxCount) {
		return false, count, nil
	}

	return true, count, nil
}
