package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	Address         string
	Database        string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	AllowedOrigins  []string
}

func loadConfig() config {
	result := config{Address: ":8080", Database: "./subsidy11.db", ShutdownTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, AllowedOrigins: []string{"*"}}
	if value := os.Getenv("SUBSIDY_ADDR"); value != "" {
		result.Address = value
	}
	if value := os.Getenv("SUBSIDY_DB"); value != "" {
		result.Database = value
	}
	if value := os.Getenv("SUBSIDY_SHUTDOWN_TIMEOUT"); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			result.ShutdownTimeout = parsed
		}
	}
	if value := os.Getenv("SUBSIDY_READ_TIMEOUT"); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			result.ReadTimeout = parsed
		}
	}
	if value := os.Getenv("SUBSIDY_WRITE_TIMEOUT"); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			result.WriteTimeout = parsed
		}
	}
	if value := os.Getenv("SUBSIDY_ORIGINS"); value != "" {
		result.AllowedOrigins = splitOrigins(value)
	}
	return result
}

func splitOrigins(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}

func envBool(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func (c config) AllowsOrigin(origin string) bool {
	if len(c.AllowedOrigins) == 0 {
		return false
	}
	for _, allowed := range c.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

func (c config) ServerSettings() map[string]string {
	return map[string]string{"address": c.Address, "database": c.Database, "read_timeout": c.ReadTimeout.String(), "write_timeout": c.WriteTimeout.String()}
}
