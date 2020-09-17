package configs

import "github.com/spf13/viper"

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
	return c.limiterHost
}
