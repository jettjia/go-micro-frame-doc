package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go-micro-frame-doc/07-consul-discovery/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("07-consul-discovery/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("07-consul-discovery/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("所有配置信息: &v", global.ServerConfig)
	fmt.Printf("%V", v.Get("name"))


}