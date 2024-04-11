package oconfig

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	TagProtocol = "protocol"
	TagHost     = "host"
)

var (
	protocol = "tcp"
)

func ParseConfig(path string) {
	viper.SetConfigName("config") // 配置文件的名字（没有扩展名）
	viper.SetConfigType("toml")   // 或者 viper.SetConfigType("TOML")
	viper.AddConfigPath(path)     // 查找配置文件的路径

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file, path: %s, err: %s \n", path, err.Error()))
	}

	if p := viper.GetString(TagProtocol); p != "" {
		protocol = p
	}
}

// Protocol 获取客户端跟服务器的通信协议
func Protocol() string {
	return protocol
}
