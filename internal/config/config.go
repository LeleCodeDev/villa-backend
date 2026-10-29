// Package config
package config

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"

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
	Location   *time.Location
}

var Env *Config

func LoadConfig() {
	_ = godotenv.Load()

	Env = &Config{
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
		Location:   getLocation(),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getLocation() *time.Location {
	name := getEnv("APP_TIMEZONE", "Asia/Jakarta")
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Fatalf("invalid APP_TIMEZONE %q: %v", name, err)
	}
	return loc
}
