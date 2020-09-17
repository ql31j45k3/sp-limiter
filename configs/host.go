package configs

import (
	"os"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	port := os.Getenv("PORT")

	if tools.IsEmpty(port) {
		port = viper.GetString("host.limiter")
	}

	config := &configHost{
		limiterHost: ":" + port,
	}

	return config
}

type configHost struct {
	limiterHost string
}

func (c *configHost) GetLimiterHost() string {
	return c.limiterHost
}
