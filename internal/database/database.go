package database

import (
	"fmt"
	"log"

	"github.com/lelecodedev/villa-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Env.DBHost,
		config.Env.DBUser,
		config.Env.DBPass,
		config.Env.DBName,
		config.Env.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB %v", err.Error()))
	}

	log.Println("database connected")

	return db
}
