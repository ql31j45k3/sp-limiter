package sp_limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ql31j45k3/sp-limiter/internal/modules/limiter"
	"go.uber.org/dig"

	"github.com/ql31j45k3/sp-limiter/configs"
)

// Start 提供對外呼叫做初始化或先執行的程式
// 不使用 init 是避免難以控制執行順序
func Start() {
	// 開始讀取設定檔，順序上必須為容器之前，執行容器內有需要設定檔 struct 取得參數
	configs.Start("")
	limiter.Start()

	// 初始化容器
	container := buildContainer()

	// 藉由容器調用 func，容器會依照 func 參數提供
	container.Invoke(limiter.RegisterRouter)
	container.Invoke(func(r *gin.Engine) {
		gin.SetMode(configs.ConfigGin.GetMode())

		r.Run(configs.ConfigHost.GetLimiterHost())
	})
}

func buildContainer() *dig.Container {
	container := dig.New()
	provideFunc := containerProvide{}

	// 置入容器管理的提供參數
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

// redisClient 建立 redis 連線
func (cp *containerProvide) redisClient() *redis.Client {
	var opt *redis.Options

	// heroku 運行環境連線方式，使用 redis.ParseURL
	if configs.ConfigRedis.GetIsProd() {
		var err error
		opt, err = redis.ParseURL(configs.ConfigRedis.GetURL())
		if err != nil {
			panic(err)
		}

		// 新版本 go-redis 如果沒有 Username = "", 會出現 Error : ERR wrong number of arguments for 'auth' command
		opt.Username = ""
	} else {
		// local redis 連線方式，使用帳密模式
		opt = &redis.Options{
			Addr:     configs.ConfigRedis.GetAddr(),
			Password: configs.ConfigRedis.GetPassword(),
			DB:       configs.ConfigRedis.GetDB(),
		}
	}

	return redis.NewClient(opt)
}
