package config

import "tokyn/pkg/utils"

type Config struct {
	SQLiteDB  string
	RedisAddr string
	RedisPass string
	AppAddr   string
}

func NewConfig() *Config {
	return &Config{
		SQLiteDB:  utils.GetEnv("SQLITE_DB", "data.db"),
		RedisAddr: utils.GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: utils.GetEnv("REDIS_PASS", ""),
		AppAddr:   utils.GetEnv("APP_ADDR", "0.0.0.0:8080"),
	}
}
