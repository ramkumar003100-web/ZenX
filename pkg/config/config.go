package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName        string        `json:"app_name" yaml:"app_name" toml:"app_name"`
	Env            string        `json:"env" yaml:"env" toml:"env"`
	HTTPAddr       string        `json:"http_addr" yaml:"http_addr" toml:"http_addr"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout" toml:"request_timeout"`
	JWTSecret      string        `json:"jwt_secret" yaml:"jwt_secret" toml:"jwt_secret"`
	DatabaseURL    string        `json:"database_url" yaml:"database_url" toml:"database_url"`
	RedisAddr      string        `json:"redis_addr" yaml:"redis_addr" toml:"redis_addr"`
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

func LoadFile(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	cfg := Load()
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(content, &cfg)
	case ".json":
		err = json.Unmarshal(content, &cfg)
	case ".toml":
		err = toml.Unmarshal(content, &cfg)
	}
	return cfg, err
}

type Watcher struct {
	path string
	mu   sync.RWMutex
	cfg  Config
}

func NewWatcher(path string) (*Watcher, error) {
	cfg, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	return &Watcher{path: path, cfg: cfg}, nil
}

func (w *Watcher) Start(interval time.Duration, onChange func(Config)) {
	go func() {
		var lastMod time.Time
		for range time.NewTicker(interval).C {
			st, err := os.Stat(w.path)
			if err != nil || !st.ModTime().After(lastMod) {
				continue
			}
			cfg, err := LoadFile(w.path)
			if err != nil {
				continue
			}
			w.mu.Lock()
			w.cfg = cfg
			w.mu.Unlock()
			lastMod = st.ModTime()
			onChange(cfg)
		}
	}()
}

func (w *Watcher) Current() Config {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.cfg
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