package internal

import (
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

	defaultHttpPort         = 8000
	defaultHttpIdleTimeout  = 30 * time.Second
	defaultHttpReadTimeout  = 10 * time.Second
	defaultHttpWriteTimeout = 10 * time.Second
)

type Config struct {
	DownloadMaxSize int
	DownloadTimeout time.Duration

	Sources      map[string]string
	SecretKey    string
	AuthToken    string
	CacheControl string
	PathPrefix   string

	HttpPort         int
	HttpIdleTimeout  time.Duration
	HttpReadTimeout  time.Duration
	HttpWriteTimeout time.Duration
}

func NewConfig() (*Config, error) {
	return &Config{
		DownloadMaxSize: getEnvInt("DOWNLOAD_MAX_SIZE", defaultDownloadMaxSize),
		DownloadTimeout: getEnvDuration("DOWNLOAD_TIMEOUT", defaultDownloadTimeout),

		Sources:      getKeyValues("SOURCES"),
		SecretKey:    getEnvString("SECRET_KEY", ""),
		AuthToken:    getEnvString("AUTH_TOKEN", ""),
		CacheControl: getEnvString("CACHE_CONTROL", defaultCacheControl),
		PathPrefix:   getEnvString("PATH_PREFIX", ""),

		HttpPort:         getEnvInt("HTTP_PORT", defaultHttpPort),
		HttpIdleTimeout:  getEnvDuration("HTTP_IDLE_TIMEOUT", defaultHttpIdleTimeout),
		HttpReadTimeout:  getEnvDuration("HTTP_READ_TIMEOUT", defaultHttpReadTimeout),
		HttpWriteTimeout: getEnvDuration("HTTP_WRITE_TIMEOUT", defaultHttpWriteTimeout),
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

func getKeyValues(key string) map[string]string {
	envVar := os.Getenv(key)
	if envVar == "" {
		panic(key + ": environment variable not set")
	}

	sources := make(map[string]string)
	pairs := strings.Split(envVar, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			panic(key + ": invalid format for environment variable, should be key1=value1;key2=value2;...")
		}
		sources[kv[0]] = kv[1]
	}
	return sources
}
