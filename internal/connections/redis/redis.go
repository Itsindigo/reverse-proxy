package redis

type RedisClient struct {
	is_connected bool
}

type RedisConnectionConfig struct{}

func (c *RedisClient) Connect() int {
	return 0
}

func (c *RedisClient) Disconnect() int {
	return 0
}
