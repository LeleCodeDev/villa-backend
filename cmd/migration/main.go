package main

import (
	"log"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/database"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func init() {
	config.LoadConfig()
}

func main() {
	db := database.NewDB()

	models := []any{
		&model.User{},
		&model.HeroSection{},
		&model.Specification{},
		&model.Faq{},
		&model.Gallery{},
		&model.SpecificationGallery{},
		&model.Application{},
		&model.Facility{},
		&model.Testimony{},
		&model.Message{},
		&model.Pricelist{},
		&model.Addon{},
		&model.AddonList{},
		&model.VillaPackage{},
		&model.VillaPackageList{},
		&model.Reason{},
		&model.RentalOption{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal("Migration failed:", err.Error())
	}

	log.Println("Successfully migrate all tables")
}
