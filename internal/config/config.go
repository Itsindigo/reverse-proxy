package config

import "fmt"

type RedisConfig struct {
	Port     string
	Host     string
	Password string
	Database string
}

type ProxyServerConfig struct {
	Port string
}

type AppConfig struct {
	Redis       RedisConfig
	ProxyServer ProxyServerConfig
}

func NewConfig() *AppConfig {
	// TODO READ FROM ENV VARS.
	return &AppConfig{
		Redis: RedisConfig{
			Port:     "6379",
			Host:     "localhost",
			Password: "",
			Database: "reverse-proxy",
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
