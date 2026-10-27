package main

import (
	"flag"
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
	var isSeedAll, isSeedAdmin, isSeedHeroSection, isSeedApplicationSetting bool

	flag.BoolVar(&isSeedAll, "all", false, "seed all data")
	flag.BoolVar(&isSeedAdmin, "admin", false, "seed admin data")
	flag.BoolVar(&isSeedHeroSection, "hero", false, "seed hero section data")
	flag.BoolVar(&isSeedApplicationSetting, "app", false, "seed application setting data")
	flag.Parse()

	db := database.NewDB()

	if isSeedAll {
		isSeedAdmin = true
		isSeedHeroSection = true
		isSeedApplicationSetting = true
	}

	if isSeedAdmin {
		seedAdmin(db)
	}
	if isSeedHeroSection {
		seedHeroSection(db)
	}
	if isSeedApplicationSetting {
		seedApplicationSetting(db)
	}

	log.Println("Successfully seed all data!")
}

func seedAdmin(db *gorm.DB) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.Env.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to seed admin")
	}

	admin := &model.User{
		Username: "Admin",
		Email:    config.Env.AdminEmail,
		Password: string(hashedPassword),
		Phone:    "00000000000",
		Role:     model.RoleAdmin,
	}

	if err := db.FirstOrCreate(admin, model.User{Email: admin.Email}).Error; err != nil {
		panic("Failed to create admin")
	}
}

func seedHeroSection(db *gorm.DB) {
	heroSection := &model.HeroSection{
		Title:    "Villa Arira",
		Subtitle: "Tempat peristirahatan eksklusif berarsitektur kayu tradisional Sunda, kolam teratai nan syahdu, private pool, meja billiard, dan halaman lapang berhawa sejuk pegunungan Lembang, Bandung Barat.",
	}

	if err := db.FirstOrCreate(heroSection, *heroSection).Error; err != nil {
		panic("Failed to create hero section")
	}
}

func seedApplicationSetting(db *gorm.DB) {
	app := &model.Application{
		Title:       "Villa Arira",
		PhoneNumber: "000000000000",
		Logo:        nil,
	}

	if err := db.FirstOrCreate(app, model.Application{Title: "Villa Arira"}).Error; err != nil {
		panic("Failed to create application setting")
	}
}
