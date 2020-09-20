package configs

import (
	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"os"

	"github.com/spf13/viper"
)

func newConfigRedis() *configRedis {
	addr := os.Getenv("REDIS_URL")
	password := ""

	if tools.IsEmpty(addr) {
		addr = viper.GetString("redis.addr")
		password = viper.GetString("redis.password")
	}


	config := &configRedis{
		addr: addr,
		password: password,
		db: viper.GetInt("redis.db"),
	}

	return config
}

type configRedis struct {
	addr string
	password string

	db int
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
