package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RDB redis.UniversalClient
)

func InitRedis() {
	cfg := GetConfig()

	if cfg.Redis.Mode == "cluster" {
		// 集群模式
		RDB = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    cfg.Redis.Cluster.Addrs,
			Password: cfg.Redis.Cluster.Password,
			ReadOnly: cfg.Redis.Cluster.ReadOnly,
		})
	} else {
		// 单机模式
		RDB = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Single.Host, cfg.Redis.Single.Port),
			Password: cfg.Redis.Single.Password,
			DB:       cfg.Redis.Single.DB,
		})
	}

	// 测试连接
	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}
}
