package limiter

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

// newRedisCounterLimiter 初始化參數
// 建立腳本的字串初始化，尚未置入 redis 內
func newRedisCounterLimiter(interval, maxCount int) *redisCounterLimiter {
	// redis script，使用 redis incr 做加總、expire 設定過期時間
	// 使用 redis script 將多個步驟達到原子性效果
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
	evalSha          string

	interval int
	maxCount int
}

// InitScriptToRedis 將腳本字串置入 redis 內做初始化
func (l *redisCounterLimiter) InitScriptToRedis(ctx context.Context, rdb *redis.Client) error {
	evalSha, err := rdb.ScriptLoad(ctx, l.counterLuaScript).Result()
	if err != nil {
		return err
	}

	l.evalSha = evalSha
	return nil
}

// TakeAvailableAndIncr 確認是否可用 (尚未到達限流條件)，並同時累加 1，代表佔一個使用量
func (l *redisCounterLimiter) TakeAvailableAndIncr(ctx context.Context, rdb *redis.Client, ip string) (bool, int64, error) {
	res, err2 := rdb.EvalSha(ctx, l.evalSha, []string{ip}, l.interval).Result()
	if err2 != nil {
		return false, 0, err2
	}

	// 將 redis 回傳做斷言，並確認是否斷言成功並免造成斷言失敗觸發 panic
	var count int64
	if v, ok := res.(int64); !ok {
		return false, 0, errors.New("redis result count Type Assertion int64 fail")
	} else {
		count = v
	}

	// 驗證是否超過 maxCount
	if count > int64(l.maxCount) {
		return false, count, nil
	}

	return true, count, nil
}
