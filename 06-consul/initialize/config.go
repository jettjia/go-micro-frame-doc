package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go-micro-frame-doc/06-consul/global"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

// 读取环境变量的配置
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("06-consul/%s-prod.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("06-consul/%s-dev.yaml", configFilePrefix)
	}

	// 读取文件配置内容
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 把内容设置到全局变量的 ServerConfig中
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	//zap.S().Infof("配置信息: &v", global.ServerConfig)

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置信息发生了变化: &v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
	})
}