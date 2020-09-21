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
	// ctx 控制多個 goroutine 的同步信號，Background 建立第一個起始點
	ctx := context.Background()
	if err := redisCounter.InitScriptToRedis(ctx, rdb); err != nil {
		panic(err)
	}

	// 建立 limiter 實作不同限流 Client 版本
	limiterRouter := newLimiter(rdb)

	// 建立 map 可用 key 取得不同限流實作版本
	apiFunc := make(map[string]func(*gin.Context))
	apiFunc[configs.HostModeCounter] = limiterRouter.getCountLimiter
	apiFunc[configs.HostModeTokenBucket] = limiterRouter.getTokenBucket
	apiFunc[configs.HostModeRedisCounter] = limiterRouter.getRedisCounter

	// 依照 ConfigHost.GetMode 參數取得限流功能
	limiterFunc, ok := apiFunc[configs.ConfigHost.GetMode()]
	if !ok {
		panic(errors.New("host.mode not exist implement func, check config [host.mode]"))
	}

	// 註冊 limiterFunc
	r.GET("/", limiterFunc)
}

func newLimiter(rdb *redis.Client) limiterRouter {
	return limiterRouter{
		rdb: rdb,
	}
}

type limiterRouter struct {
	rdb *redis.Client
}

func (l *limiterRouter) getCountLimiter(c *gin.Context) {
	// 取得 客戶端 IP
	clientIP := c.ClientIP()

	if countLimit.TakeAvailableAndIncr(clientIP) {
		// 未達成限流條件，回傳目前的請求量
		c.String(http.StatusOK, countLimit.GetCount(clientIP))
		return
	}

	// 達成限流條件，回傳 Error
	c.String(http.StatusOK, "Error")
}

func (l *limiterRouter) getTokenBucket(c *gin.Context) {
	// 取得 客戶端 IP
	clientIP := c.ClientIP()
	// 不阻塞方式取得 Token
	block := false

	if ok, count := tokenBucket.TakeAvailable(clientIP, block); ok {
		// 未達成限流條件，回傳目前的請求量
		c.String(http.StatusOK, strconv.Itoa(int(count)))
		return
	}

	// 達成限流條件，回傳 Error
	c.String(http.StatusOK, "Error")
}

func (l *limiterRouter) getRedisCounter(c *gin.Context) {
	// ctx 控制多個 goroutine 的同步信號，Background 建立第一個起始點
	ctx := context.Background()
	// 取得 客戶端 IP
	clientIP := c.ClientIP()

	ok, count, err := redisCounter.TakeAvailableAndIncr(ctx, l.rdb, clientIP)
	if err != nil {
		log.Println(fmt.Sprintf("Error TakeAvailableAndIncr fail %v", err))
		c.String(http.StatusOK, "Error TakeAvailableAndIncr fail")
		return
	}

	// 未達成限流條件，回傳目前的請求量
	if ok {
		c.String(http.StatusOK, strconv.Itoa(int(count)))
		return
	}

	// 達成限流條件，回傳 Error
	c.String(http.StatusOK, "Error")
}
