package bootstrap

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置文件
func InitConfig(file string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("找不到配置文件:%s", file)) // 系统初始化阶段发生任何错误，直接结束进程
		}
		panic(fmt.Errorf("解析配置文件%s出错:%s", file, err))
	}
	v.WatchConfig()
	return v
}
