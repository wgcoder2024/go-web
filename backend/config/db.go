package config

import (
	"log"

	"github.com/wgcoder2024/go-web/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	cfg := GetConfig()
	DB, err = gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移
	DB.AutoMigrate(&models.User{})
}
