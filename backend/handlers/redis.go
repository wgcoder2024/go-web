package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wgcoder2024/go-web/backend/config"
)

type RedisEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
	TTL   int64  `json:"ttl"`
	Node  string `json:"node,omitempty"`
}

// GetRedisKeys 获取所有键
func GetRedisKeys(c *gin.Context) {
	pattern := c.DefaultQuery("pattern", "*")
	ctx := context.Background()

	var entries []RedisEntry

	// 获取所有键
	var keys []string
	var err error

	if cluster, ok := config.RDB.(*redis.ClusterClient); ok {
		// 集群模式：遍历所有节点获取键
		err = cluster.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			nodeKeys, err := client.Keys(ctx, pattern).Result()
			if err != nil {
				return err
			}

			for _, key := range nodeKeys {
				keyType, _ := client.Type(ctx, key).Result()
				value, _ := client.Get(ctx, key).Result()
				ttl := config.RDB.TTL(ctx, key).Val().Seconds()
				node := client.Options().Addr

				entries = append(entries, RedisEntry{
					Key:   key,
					Value: value,
					Type:  keyType,
					TTL:   int64(ttl),
					Node:  node,
				})
			}
			return nil
		})
	} else {
		// 单机模式
		keys, err = config.RDB.Keys(ctx, pattern).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, key := range keys {
			keyType, _ := config.RDB.Type(ctx, key).Result()
			value, _ := config.RDB.Get(ctx, key).Result()
			ttl := config.RDB.TTL(ctx, key).Val().Seconds()

			entries = append(entries, RedisEntry{
				Key:   key,
				Value: value,
				Type:  keyType,
				TTL:   int64(ttl),
			})
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// SetRedisKey 设置键值
func SetRedisKey(c *gin.Context) {
	var entry RedisEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	var err error
	if entry.TTL > 0 {
		err = config.RDB.Set(ctx, entry.Key, entry.Value, time.Duration(entry.TTL)*time.Second).Err()
	} else {
		err = config.RDB.Set(ctx, entry.Key, entry.Value, 0).Err()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "键值对设置成功"})
}

// DeleteRedisKey 删除键
func DeleteRedisKey(c *gin.Context) {
	key := c.Param("key")
	ctx := context.Background()

	err := config.RDB.Del(ctx, key).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "键删除成功"})
}
