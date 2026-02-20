package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName        string
	Env            string
	HTTPAddr       string
	RequestTimeout time.Duration
	JWTSecret      string
	DatabaseURL    string
	RedisAddr      string
}

func Load() Config {
	return Config{
		AppName:        get("ZENX_APP_NAME", "ZenX"),
		Env:            get("ZENX_ENV", "development"),
		HTTPAddr:       get("ZENX_HTTP_ADDR", ":8080"),
		RequestTimeout: time.Duration(getInt("ZENX_REQUEST_TIMEOUT_SEC", 15)) * time.Second,
		JWTSecret:      get("ZENX_JWT_SECRET", "dev-secret"),
		DatabaseURL:    get("ZENX_DATABASE_URL", ""),
		RedisAddr:      get("ZENX_REDIS_ADDR", "127.0.0.1:6379"),
	}
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return fallback
}
