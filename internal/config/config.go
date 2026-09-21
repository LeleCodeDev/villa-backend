package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	BaseURL   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

var Env *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Failed to load environment %v", err.Error()))
	}

	config := &Config{
		Port:      getEnv("PORT", "3000"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "villa_aira"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:3000"),
		JWTSecret: getEnv("JWT_SECRET", "SECRET"),
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
