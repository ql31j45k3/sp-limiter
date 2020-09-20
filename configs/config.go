package configs

import (
	"github.com/ql31j45k3/sp-limiter/internal/utils/tools"
	"os"

	"github.com/spf13/viper"
)

var (
	ConfigHost  *configHost
	ConfigGin   *configGin
	ConfigRedis *configRedis
)

// Start 開始 Config 設定參數與讀取檔案並轉成 struct
// 預設會抓取執行程式的啟示點資料夾，可用參數調整路徑來源
func Start(sourcePath string) {
	viper.AddConfigPath(getPath(sourcePath))
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	ConfigHost = newConfigHost()
	ConfigGin = newConfigGin()
	ConfigRedis = newConfigRedis()
}

// getPath 預設會抓取執行程式的啟示點資料夾
// 可用參數調整路徑來源
func getPath(sourcePath string) string {
	if tools.IsNotEmpty(sourcePath) {
		return sourcePath + "/configs"
	}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path + "/configs"
}
