package configs

import (
	"os"

	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"github.com/spf13/viper"
)

func newConfigHost() *configHost {
	config := &configHost{
		limiterHost: ":" + viper.GetString("host.limiter"),
	}

	return config
}

type configHost struct {
	limiterHost string
}

func (c *configHost) GetLimiterHost() string {
	port := os.Getenv("PORT")
	if tools.IsEmpty(port) {
		return c.limiterHost
	}

	return port
}
