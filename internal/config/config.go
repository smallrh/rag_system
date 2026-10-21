package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Load 加载应用配置
func Load() (*Config, error) {
	// .env 不存在时，允许直接使用系统环境变量
	_ = godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", ""),
			SecretKey: getEnv("MINIO_SECRET_KEY", ""),
			Bucket:    getEnv("MINIO_BUCKET", "auth-files"),
			UseSSL:    getEnvBool("MINIO_USE_SSL", false),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 校验必要配置
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL 未配置")
	}

	if c.MinIO.Endpoint == "" {
		return fmt.Errorf("MINIO_ENDPOINT 未配置")
	}

	if c.MinIO.AccessKey == "" {
		return fmt.Errorf("MINIO_ACCESS_KEY 未配置")
	}

	if c.MinIO.SecretKey == "" {
		return fmt.Errorf("MINIO_SECRET_KEY 未配置")
	}

	if c.MinIO.Bucket == "" {
		return fmt.Errorf("MINIO_BUCKET 未配置")
	}

	if c.Redis.DB < 0 {
		return fmt.Errorf("REDIS_DB 不能小于 0")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return n
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}
