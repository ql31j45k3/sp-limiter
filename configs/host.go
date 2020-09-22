package configs

import (
	"os"
	"strings"
	"time"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"github.com/spf13/viper"
)

const (
	HostModeCounter      = "counter"
	HostModeTokenBucket  = "tokenbucket"
	HostModeRedisCounter = "rediscounter"
)

func newConfigHost() *configHost {
	// 取得 heroku 運行環境的 PORT 變數
	port := os.Getenv("PORT")
	if tools.IsEmpty(port) {
		port = viper.GetString("host.limiterPort")
	}

	viper.SetDefault("host.mode", HostModeCounter)
	viper.SetDefault("host.interval", 60)
	viper.SetDefault("host.maxCount", 60)

	config := &configHost{
		limiterHost: ":" + port,
		mode:        viper.GetString("host.mode"),
		interval:    viper.GetDuration("host.interval"),
		intervalInt: viper.GetInt("host.interval"),
		maxCount:    viper.GetInt("host.maxCount"),
	}

	return config
}

type configHost struct {
	limiterHost string

	mode        string
	interval    time.Duration
	intervalInt int
	maxCount    int
}

func (c *configHost) GetLimiterHost() string {
	return c.limiterHost
}

func (c *configHost) GetMode() string {
	return strings.ToLower(c.mode)
}

func (c *configHost) GetInterval() time.Duration {
	return c.interval * time.Second
}

func (c *configHost) GetIntervalInt() int {
	return c.intervalInt
}

func (c *configHost) GetMaxCount() int {
	return c.maxCount
}
