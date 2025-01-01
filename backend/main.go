package main

import (
	"fmt"

	"github.com/wgcoder2024/go-web/backend/config"
	"github.com/wgcoder2024/go-web/backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	// 初始化数据库
	config.InitDB()

	// 创建 Gin 路由
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.Cors.AllowedOrigins,
		AllowMethods:     cfg.Server.Cors.AllowedMethods,
		AllowHeaders:     cfg.Server.Cors.AllowedHeaders,
		AllowCredentials: true,
	}))

	// API 路由
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
		}

		// 数据库表管理路由
		tables := v1.Group("/tables")
		{
			tables.GET("", handlers.GetTables)
			tables.GET("/:name", handlers.GetTableDetails)
			tables.POST("", handlers.CreateTable)
			tables.DELETE("/:name", handlers.DeleteTable)

			// 新增路由
			tables.PUT("/:name", handlers.AlterTable)
			tables.GET("/:name/data", handlers.GetTableData)
			tables.POST("/query", handlers.ExecuteSQL)
			tables.GET("/:name/export", handlers.ExportTable)
			tables.POST("/import", handlers.ImportTable)
		}

		redis := v1.Group("/redis")
		{
			redis.GET("/keys", handlers.GetRedisKeys)
			redis.POST("/keys", handlers.SetRedisKey)
			redis.DELETE("/keys/:key", handlers.DeleteRedisKey)
		}
	}

	// 启动服务器
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	r.Run(port)
}
