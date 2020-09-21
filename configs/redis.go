package configs

import (
	"os"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"github.com/spf13/viper"
)

func newConfigRedis() *configRedis {
	isProd := false
	// 取得 heroku 運行環境的 REDIS_URL 變數
	url := os.Getenv("REDIS_URL")

	if tools.IsNotEmpty(url) {
		isProd = true
	}

	config := &configRedis{
		isProd:   isProd,
		url:      url,
		addr:     viper.GetString("redis.addr"),
		password: viper.GetString("redis.password"),
		db:       viper.GetInt("redis.db"),
	}

	return config
}

type configRedis struct {
	isProd bool

	url      string
	addr     string
	password string

	db int
}

func (c *configRedis) GetIsProd() bool {
	return c.isProd
}

func (c *configRedis) GetURL() string {
	return c.url
}

func (c *configRedis) GetAddr() string {
	return c.addr
}

func (c *configRedis) GetPassword() string {
	return c.password
}

func (c *configRedis) GetDB() int {
	return c.db
}
