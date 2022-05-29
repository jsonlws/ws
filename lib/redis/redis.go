package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

//初始化reids连接对象
func NewRedis(db int) *redis.Client {
	RedisObj := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.pwd"),    // redis连接密码
		DB:       db,                              // 使用redis的库
		PoolSize: viper.GetInt("redis.pool_size"), // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := RedisObj.Ping(ctx).Result()
	if err != nil {
		panic(pong)
	}
	return RedisObj
}
