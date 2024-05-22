package app_config

import (
	"fmt"
)

type RedisConfig struct {
	Port     string
	Host     string
	Password string
	Database int
}

type ProxyServerConfig struct {
	Port string
}

type AppConfig struct {
	Redis       RedisConfig
	ProxyServer ProxyServerConfig
}

func NewConfig() *AppConfig {
	return &AppConfig{
		Redis: RedisConfig{
			Port:     "6379",
			Host:     "redis",
			Password: "boop",
			Database: 0,
		},
		ProxyServer: ProxyServerConfig{
			Port: "6666",
		},
	}
}

func (c AppConfig) String() string {
	return fmt.Sprintf(
		"AppConfig(Redis(Port: %q, Host: %q, Password: %q, Database %q), ProxyServer())",
		c.Redis.Port,
		c.Redis.Host,
		c.Redis.Password,
		c.Redis.Database,
	)
}
