package oconfig

import (
	"github.com/spf13/viper"
)

const (
	TagProtocol = "protocol"
	TagIp       = "ip"
	TagPort     = "port"
)

var (
	protocol = "tcp"
)

func ParseConfig(path string) {
	viper.SetConfigType("toml") // 或者 viper.SetConfigType("TOML")
	viper.SetConfigFile(path)   // 查找配置文件的路径

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if p := viper.GetString(TagProtocol); p != "" {
		protocol = p
	}

	if tp := viper.GetString(TagPort); tp == "" {
		viper.Set(TagPort, "8900")
	}
}

// Protocol 获取客户端跟服务器的通信协议
func Protocol() string {
	return protocol
}
