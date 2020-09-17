package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/ql31j45k3/sp-limiter/configs"
)

var (
	countLimit *counterLimit

	apiFunc map[string]func(*gin.Context)
)

func Start() {
	countLimit = newCounterLimit(configs.ConfigHost.GetInterval(), configs.ConfigHost.GetMaxCount())

	limiterRouter := newLimiter()

	apiFunc = make(map[string]func(*gin.Context))

	apiFunc[configs.HostModeCounter] = limiterRouter.getCountLimiter
}
