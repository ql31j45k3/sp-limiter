package limiter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ql31j45k3/sp-limiter/configs"
)

func RegisterRouter(r *gin.Engine, rdb *redis.Client) {
	limiterRouter := newLimiter(rdb)

	apiFunc := make(map[string]func(*gin.Context))
	apiFunc[configs.HostModeCounter] = limiterRouter.getCountLimiter
	apiFunc[configs.HostModeTokenBucket] = limiterRouter.getTokenBucket
	apiFunc[configs.HostModeRedis] = limiterRouter.getRedisLimiter

	r.GET("/", apiFunc[configs.ConfigHost.GetMode()])
}

func newLimiter(rdb *redis.Client) limiterRouter {
	return limiterRouter{
		rdb:rdb,
	}
}

type limiterRouter struct {
	rdb *redis.Client
}

func (l *limiterRouter) getCountLimiter(c *gin.Context) {
	clientIP := c.ClientIP()

	if countLimit.IsAvailable(clientIP) {
		countLimit.Increase(clientIP)
		c.String(http.StatusOK, countLimit.GetCount(clientIP))
		return
	}

	c.String(http.StatusOK, "Error")
}

func (l *limiterRouter) getTokenBucket(c *gin.Context) {
	clientIP := c.ClientIP()
	block := false

	if ok, count := tokenBucket.TakeAvailable(clientIP, block); ok {
		c.String(http.StatusOK, strconv.Itoa(int(count)))
		return
	}

	c.String(http.StatusOK, "Error")
}

func (l *limiterRouter) getRedisLimiter(c *gin.Context) {
	ctx := context.Background()

	v, err := l.rdb.Get(ctx, "test").Result()
	fmt.Println(v, err)

	c.String(http.StatusOK, v)
}
