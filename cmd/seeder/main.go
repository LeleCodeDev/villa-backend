package main

import (
	"log"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/database"
	"github.com/lelecodedev/villa-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	config.LoadConfig()
}

func main() {
	db := database.NewDB()
	seedAdmin(db)

	log.Println("Successfully seed all data!")
}

func seedAdmin(db *gorm.DB) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin#1234"), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to seed admin")
	}

	admin := &model.User{
		Username: "Admin",
		Email:    "admin@email.xyz",
		Password: string(hashedPassword),
		Phone:    "00000000000",
		Role:     model.RoleAdmin,
	}

	if err := db.FirstOrCreate(admin, model.User{Email: admin.Email}).Error; err != nil {
		panic("Failed to create admin")
	}
}
