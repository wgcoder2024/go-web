package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
		Cors struct {
			AllowedOrigins []string `yaml:"allowed_origins"`
			AllowedMethods []string `yaml:"allowed_methods"`
			AllowedHeaders []string `yaml:"allowed_headers"`
		} `yaml:"cors"`
	} `yaml:"server"`

	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Params   string `yaml:"params"`
	} `yaml:"database"`

	mu sync.RWMutex
}

var (
	cfg  *Config
	once sync.Once
)

// GetConfig 返回配置单例
func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := cfg.Load(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		go cfg.Watch()
	})
	return cfg
}

// Load 加载配置文件
func (c *Config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// Watch 监听配置文件变化
func (c *Config) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create watcher: %v", err)
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file modified, reloading...")
					time.Sleep(100 * time.Millisecond) // 等待文件写入完成
					if err := c.Load(); err != nil {
						log.Printf("Failed to reload config: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	if err := watcher.Add("config/config.yaml"); err != nil {
		log.Printf("Failed to watch config file: %v", err)
	}
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.Params,
	)
}
