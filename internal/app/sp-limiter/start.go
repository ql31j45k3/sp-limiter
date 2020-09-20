package sp_limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ql31j45k3/sp-limiter/internal/modules/limiter"
	"go.uber.org/dig"

	"github.com/ql31j45k3/sp-limiter/configs"
)

func Start() {
	configs.Start("")
	limiter.Start()

	container := buildContainer()

	container.Invoke(limiter.RegisterRouter)
	container.Invoke(func(r *gin.Engine) {
		gin.SetMode(configs.ConfigGin.GetMode())

		r.Run(configs.ConfigHost.GetLimiterHost())
	})
}

func buildContainer() *dig.Container {
	container := dig.New()
	provideFunc := containerProvide{}

	container.Provide(provideFunc.gin)
	container.Provide(provideFunc.redisClient)

	return container
}

type containerProvide struct {
}

// gin 建立 gin Engine，設定 middleware
func (cp *containerProvide) gin() *gin.Engine {
	return gin.Default()
}

func (cp *containerProvide) redisClient() *redis.Client {
	var opt *redis.Options

	if configs.ConfigRedis.GetIsProd() {
		var err error
		opt, err = redis.ParseURL(configs.ConfigRedis.GetURL())
		if err != nil {
			panic(err)
		}

		opt.Username = ""

	} else {
		opt = &redis.Options{
			Addr:     configs.ConfigRedis.GetAddr(),
			Password: configs.ConfigRedis.GetPassword(),
			DB:       configs.ConfigRedis.GetDB(),
		}
	}

	return redis.NewClient(opt)
}
