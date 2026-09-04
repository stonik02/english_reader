package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment             string
	GRPCAddr                string
	HTTPAddr                string
	FrontendOrigin          string
	DatabaseURL             string
	JWTSecret               string
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	ShutdownTimeout         time.Duration
	BookStoragePath         string
	TranslateURL            string
	TranslateTimeout        time.Duration
	TranslateCacheTTL       time.Duration
	LookupWordMaxLength     int
	LookupSentenceMaxLength int
}

func Load() (Config, error) {
	shutdownTimeout, err := time.ParseDuration(value("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SHUTDOWN_TIMEOUT: %w", err)
	}
	accessTokenTTL, err := time.ParseDuration(value("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		return Config{}, fmt.Errorf("parse ACCESS_TOKEN_TTL: %w", err)
	}
	refreshTokenTTL, err := time.ParseDuration(value("REFRESH_TOKEN_TTL", "720h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse REFRESH_TOKEN_TTL: %w", err)
	}
	translateTimeout, err := time.ParseDuration(value("TRANSLATE_TIMEOUT", "3s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse TRANSLATE_TIMEOUT: %w", err)
	}
	translateCacheTTL, err := time.ParseDuration(value("TRANSLATE_CACHE_TTL", "168h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse TRANSLATE_CACHE_TTL: %w", err)
	}

	cfg := Config{
		Environment:             value("APP_ENV", "development"),
		GRPCAddr:                value("GRPC_ADDR", ":9090"),
		HTTPAddr:                value("HTTP_ADDR", ":8081"),
		FrontendOrigin:          value("FRONTEND_ORIGIN", "http://localhost:5173"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		JWTSecret:               os.Getenv("JWT_SECRET"),
		AccessTokenTTL:          accessTokenTTL,
		RefreshTokenTTL:         refreshTokenTTL,
		ShutdownTimeout:         shutdownTimeout,
		BookStoragePath:         value("BOOK_STORAGE_PATH", "/var/lib/english-reader/books"),
		TranslateURL:            value("TRANSLATE_URL", "http://localhost:5000"),
		TranslateTimeout:        translateTimeout,
		TranslateCacheTTL:       translateCacheTTL,
		LookupWordMaxLength:     intValue("LOOKUP_WORD_MAX_LENGTH", 96),
		LookupSentenceMaxLength: intValue("LOOKUP_SENTENCE_MAX_LENGTH", 1000),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if cfg.AccessTokenTTL <= 0 || cfg.RefreshTokenTTL <= 0 {
		return Config{}, fmt.Errorf("token TTLs must be positive")
	}
	if cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be positive")
	}
	if cfg.TranslateTimeout <= 0 || cfg.TranslateCacheTTL <= 0 || cfg.LookupWordMaxLength <= 0 || cfg.LookupSentenceMaxLength <= 0 {
		return Config{}, fmt.Errorf("translator configuration values must be positive")
	}

	return cfg, nil
}

func intValue(key string, fallback int) int {
	parsed, err := strconv.Atoi(value(key, strconv.Itoa(fallback)))
	if err != nil {
		return fallback
	}
	return parsed
}

func value(key, fallback string) string {
	if current, ok := os.LookupEnv(key); ok && current != "" {
		return current
	}
	return fallback
}
