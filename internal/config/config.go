package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	MySQL    MySQLConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Mailtrap MailtrapConfig
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Addr:   fmt.Sprintf(":%s", getEnv("SERVER_PORT", "8080")),
			APIKey: getEnv("SERVER_API_KEY", "MY_API_KEY"),
		},
		MySQL: MySQLConfig{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "3306"),
			DBUser:     getEnv("DB_USER", "root"),
			DBPassword: getEnv("DB_PASSWORD", "123456"),
			DBName:     getEnv("DB_NAME", "quiz"),
		},
		JWT: JWTConfig{
			SecretKey:       []byte(getEnv("JWT_SECRET_KEY", "MY_SECRET_KEY")),
			AccessTokenTTL:  getEnvAsDuration("ACCESS_TOKEN_TTL", 2*time.Hour),
			RefreshTokenTTL: getEnvAsDuration("REFRESH_TOKEN_TTL", 72*time.Hour),
		},
		Redis: RedisConfig{
			Addr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
			Password: getEnv("REDIS_PASSWORD", "123456"),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Mailtrap: MailtrapConfig{
			Host:     getEnv("MAILTRAP_HOST", ""),
			Port:     getEnvAsInt("MAILTRAP_PORT", 0),
			Username: getEnv("MAILTRAP_USERNAME", ""),
			Password: getEnv("MAILTRAP_PASSWORD", ""),
		},
	}
}

type ServerConfig struct {
	Addr   string
	APIKey string
}

type MySQLConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type JWTConfig struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type MailtrapConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}

	return fallback
}
