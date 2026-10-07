// Package config
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	BaseURL    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	JWTSecret  string
	AdminEmail string
	AdminPass  string
}

var Env *Config

func LoadConfig() {
	_ = godotenv.Load()

	config := &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", ""),
		DBName:     getEnv("DB_NAME", "villa_aira"),
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080"),
		JWTSecret:  getEnv("JWT_SECRET", "SECRET"),
		AdminEmail: getEnv("ADMIN_EMAIL", "admin@admin.com"),
		AdminPass:  getEnv("ADMIN_PASS", ""),
	}

	Env = config
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}

	return value
}
