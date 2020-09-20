package limiter

import (
	"github.com/ql31j45k3/sp-limiter/configs"
)

var (
	countLimit *counterLimit
	tokenBucket *tokenBucketLimiter
)

func Start() {
	countLimit = newCounterLimit(configs.ConfigHost.GetInterval(), configs.ConfigHost.GetMaxCount())
	tokenBucket = newTokenBucketLimiter(configs.ConfigHost.GetInterval(), int64(configs.ConfigHost.GetMaxCount()))
}
