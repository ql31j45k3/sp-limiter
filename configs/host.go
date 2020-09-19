package configs

import (
	"os"
	"strings"
	"time"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"github.com/spf13/viper"
)

const (
	HostModeCounter = "counter"
	HostModeTokenBucket = "tokenbucket"
)

func newConfigHost() *configHost {
	port := os.Getenv("PORT")
	if tools.IsEmpty(port) {
		port = viper.GetString("host.limiter")
	}

	viper.SetDefault("host.mode", HostModeCounter)
	viper.SetDefault("host.interval", 60)
	viper.SetDefault("host.maxCount", 60)

	config := &configHost{
		limiterHost: ":" + port,
		mode: viper.GetString("host.mode"),
		interval: viper.GetDuration("host.interval"),
		maxCount: viper.GetInt("host.maxCount"),
	}

	return config
}

type configHost struct {
	limiterHost string

	mode string
	interval time.Duration
	maxCount int
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

func (c *configHost) GetMaxCount() int {
	return c.maxCount
}
