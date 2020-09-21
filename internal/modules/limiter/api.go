package limiter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ql31j45k3/sp-limiter/configs"
)

func RegisterRouter(r *gin.Engine, rdb *redis.Client) {
	ctx := context.Background()
	if err := redisCounter.initScript(ctx, rdb); err != nil {
		panic(err)
	}

	limiterRouter := newLimiter(rdb)

	apiFunc := make(map[string]func(*gin.Context))
	apiFunc[configs.HostModeCounter] = limiterRouter.getCountLimiter
	apiFunc[configs.HostModeTokenBucket] = limiterRouter.getTokenBucket
	apiFunc[configs.HostModeRedisCounter] = limiterRouter.getRedisCounter

	limiterFunc, ok := apiFunc[configs.ConfigHost.GetMode()];
	if !ok {
		panic(errors.New("host.mode not exist implement func"))
	}

	r.GET("/", limiterFunc)
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

	if countLimit.TakeAvailableAndIncr(clientIP) {
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

func (l *limiterRouter) getRedisCounter(c *gin.Context) {
	ctx := context.Background()
	clientIP := c.ClientIP()

	ok, count, err := redisCounter.TakeAvailableAndIncr(ctx, l.rdb, clientIP)
	if err != nil {
		log.Println(fmt.Sprintf("Error TakeAvailableAndIncr fail %v", err))
		c.String(http.StatusOK, "Error TakeAvailableAndIncr fail")
		return
	}

	if ok {
		c.String(http.StatusOK, strconv.Itoa(int(count)))
		return
	}

	c.String(http.StatusOK, "Error")
}
