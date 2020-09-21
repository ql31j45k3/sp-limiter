package configs

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	ginModeDebug = "debug"
)

func newConfigGin() *configGin {
	viper.SetDefault("gin.mode", ginModeDebug)

	config := &configGin{
		mode: viper.GetString("gin.mode"),
	}

	return config
}

type configGin struct {
	mode string
}

func (c *configGin) GetMode() string {
	return strings.ToLower(c.mode)
}
