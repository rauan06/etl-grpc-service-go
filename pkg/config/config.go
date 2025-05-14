package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Config struct {
	DBHost        string `env:"DB_HOST"`
	DBUser        string `env:"DB_USER"`
	DBPassword    string `env:"DB_PASSWORD"`
	DBName        string `env:"DB_NAME"`
	DBPort        string `env:"DB_PORT"`
	JWTSecret     string `env:"JWT_SECRET"`
	RedisURI      string `env:"REDIS_URI"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
}

var cfg Config

func LoadConfig() *Config {
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "postgres")
	cfg.DBName = getEnv("DB_NAME", "cafeteria")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.JWTSecret = createMd5Hash(getEnv("JWT_SECRET", "not-so-secret-now-is-it?"))
	cfg.RedisURI = getEnv("REDIS_URI", "redis:6379")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB, _ = strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &cfg
}

func GetConfing() *Config {
	return &cfg
}

func GetJWTSecret() string {
	return cfg.JWTSecret
}

func (c *Config) MakeConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func createMd5Hash(text string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, text)
	if err != nil {
		// panic(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
