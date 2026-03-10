package internal

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	KB = 1024
	MB = 1024 * KB

	defaultDownloadMaxSize = 50 * MB
	defaultDownloadTimeout = 10 * time.Second

	defaultCacheControl = "public, max-age=31536000"

	defaultMaxConcurrentTransforms = 10

	defaultHttpPort         = 8000
	defaultHttpIdleTimeout  = 30 * time.Second
	defaultHttpReadTimeout  = 10 * time.Second
	defaultHttpWriteTimeout = 10 * time.Second
)

type SourceConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	DownloadMaxSize int
	DownloadTimeout time.Duration

	Sources      []SourceConfig
	Renderers    []SourceConfig
	SecretKey    string
	AuthToken    string
	MaxConcurrentTransforms int
	CacheControl            string
	PathPrefix              string

	HttpPort         int
	HttpIdleTimeout  time.Duration
	HttpReadTimeout  time.Duration
	HttpWriteTimeout time.Duration
}

func NewConfig() (*Config, error) {
	sources, err := getSourceConfigs("MEDIATOR_SOURCES")
	if err != nil {
		return nil, err
	}

	renderers, err := getSourceConfigs("MEDIATOR_RENDERERS")
	if err != nil {
		return nil, err
	}

	return &Config{
		DownloadMaxSize: getEnvInt("MEDIATOR_DOWNLOAD_MAX_SIZE", defaultDownloadMaxSize),
		DownloadTimeout: getEnvDuration("MEDIATOR_DOWNLOAD_TIMEOUT", defaultDownloadTimeout),

		MaxConcurrentTransforms: getEnvInt("MEDIATOR_MAX_CONCURRENT_TRANSFORMS", defaultMaxConcurrentTransforms),

		Sources: sources,
		Renderers:    renderers,
		SecretKey:    getEnvString("MEDIATOR_SECRET_KEY", ""),
		AuthToken:    getEnvString("MEDIATOR_AUTH_TOKEN", ""),
		CacheControl: getEnvString("MEDIATOR_CACHE_CONTROL", defaultCacheControl),
		PathPrefix:   getEnvString("MEDIATOR_PATH_PREFIX", ""),

		HttpPort:         getEnvInt("MEDIATOR_HTTP_PORT", defaultHttpPort),
		HttpIdleTimeout:  getEnvDuration("MEDIATOR_HTTP_IDLE_TIMEOUT", defaultHttpIdleTimeout),
		HttpReadTimeout:  getEnvDuration("MEDIATOR_HTTP_READ_TIMEOUT", defaultHttpReadTimeout),
		HttpWriteTimeout: getEnvDuration("MEDIATOR_HTTP_WRITE_TIMEOUT", defaultHttpWriteTimeout),
	}, nil
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return time.Duration(intValue) * time.Second
}

func getSourceConfigs(key string) ([]SourceConfig, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		return []SourceConfig{}, nil
	}

	envVar = strings.TrimSpace(envVar)

	var result []SourceConfig
	if err := json.Unmarshal([]byte(envVar), &result); err != nil {
		return nil, fmt.Errorf("%s: invalid JSON format for environment variable, should be a JSON array like [{ \"name\": \"source1\", \"url\": \"http://example.com\" }]: %w", key, err)
	}

	slog.Debug("parsing config", "key", key, "parsed", fmt.Sprintf("%+v", result))

	return result, nil
}

func (c *Config) FindSourceByName(name string) (string, bool) {
	for _, source := range c.Sources {
		if source.Name == name {
			return source.URL, true
		}
	}
	return "", false
}

func (c *Config) FindRendererByName(name string) (string, bool) {
	for _, renderer := range c.Renderers {
		if renderer.Name == name {
			return renderer.URL, true
		}
	}
	return "", false
}
