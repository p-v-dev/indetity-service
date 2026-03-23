package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application-wide configuration loaded from environment variables.
// No config values are hardcoded — everything comes from the environment.
type Config struct {
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
}

// LoadConfig reads environment variables (and an optional .env file) and returns
// a populated Config. It calls log.Fatal if any required value is missing or invalid.
func LoadConfig() *Config {
	// In development, values can live in a .env file.
	// In production (Docker / Kubernetes) they are injected directly — no .env file is needed.
	if err := godotenv.Load(); err != nil {
		log.Println("config: no .env file found, reading from environment directly")
	}

	cfg := &Config{
		Port:        requireEnv("PORT"),
		DatabaseURL: requireEnv("DATABASE_URL"),
		RedisURL:    requireEnv("REDIS_URL"),
		JWTSecret:   requireEnv("JWT_SECRET"),
	}

	var err error

	cfg.JWTAccessTTL, err = parseDuration(requireEnv("JWT_ACCESS_TTL"))
	if err != nil {
		log.Fatalf("config: invalid JWT_ACCESS_TTL: %v", err)
	}

	cfg.JWTRefreshTTL, err = parseDuration(requireEnv("JWT_REFRESH_TTL"))
	if err != nil {
		log.Fatalf("config: invalid JWT_REFRESH_TTL: %v", err)
	}

	return cfg
}

// requireEnv returns the value of an environment variable or calls log.Fatal if it is unset.
func requireEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("config: required environment variable %q is not set", key)
	}
	return v
}

// parseDuration extends time.ParseDuration with support for the "d" (day) suffix.
// Examples: "15m", "1h", "7d".
//
// Go's standard library only understands up to hours ("h"), so "7d" would fail without
// this wrapper. We handle "d" manually and delegate everything else to time.ParseDuration.
func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		raw := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("parseDuration: %q is not a valid duration (expected a number before 'd')", s)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("parseDuration: %w", err)
	}
	return d, nil
}
