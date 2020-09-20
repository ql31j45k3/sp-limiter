package limiter

import (
	"github.com/ql31j45k3/sp-limiter/configs"
)

var (
	countLimit *counterLimit
	tokenBucket *tokenBucketLimiter
	redisCounter *redisCounterLimiter
)

func Start() {
	countLimit = newCounterLimit(configs.ConfigHost.GetInterval(), configs.ConfigHost.GetMaxCount())
	tokenBucket = newTokenBucketLimiter(configs.ConfigHost.GetInterval(), int64(configs.ConfigHost.GetMaxCount()))
	redisCounter = newRedisCounterLimiter(configs.ConfigHost.GetIntervalInt(), configs.ConfigHost.GetMaxCount())
}
