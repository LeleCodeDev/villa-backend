package bootstrap

import (
	"reflect"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lelecodedev/villa-backend/internal/database"
	"github.com/lelecodedev/villa-backend/internal/handler"
	"github.com/lelecodedev/villa-backend/internal/middleware"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/internal/service"
	"gorm.io/gorm"
)

type App struct {
	Router                  *gin.Engine
	DB                      *gorm.DB
	AuthMiddleware          gin.HandlerFunc
	AuthHandler             *handler.AuthHandler
	ApplicationHandler      *handler.ApplicationHandler
	HeroSectionHandler      *handler.HeroSectionHandler
	FaqHandler              *handler.FaqHandler
	GalleryHandler          *handler.GalleryHandler
	TestimonyHandler        *handler.TestimonyHandler
	VillaPackageListHandler *handler.VillaPackageListHandler
}

func NewApp() *App {
	setupValidator()

	r := gin.Default()
	db := database.NewDB()
	txManager := repository.NewTxManager(db)

	applicationRepo := repository.NewApplicationRepository(db)
	userRepo := repository.NewUserRepository(db)
	heroSectionRepo := repository.NewHeroSectionRepository(db)
	faqRepo := repository.NewFaqRepository(db)
	galleryRepo := repository.NewGalleryRepository(db)
	testimonyRepo := repository.NewTestimonyRepository(db)
	villaPackageListRepo := repository.NewVillaPackageListRepository(db)

	authService := service.NewAuthService(txManager, userRepo)
	applicationService := service.NewApplicationService(txManager, applicationRepo)
	heroSectionService := service.NewHeroSectionService(txManager, heroSectionRepo)
	faqService := service.NewFaqService(txManager, faqRepo)
	galleryService := service.NewGalleryService(txManager, galleryRepo)
	testimonyService := service.NewTestimonyService(txManager, testimonyRepo)
	villaPackageListService := service.NewVillaPackageListService(txManager, villaPackageListRepo)

	app := &App{
		Router:                  r,
		DB:                      db,
		AuthMiddleware:          middleware.AuthMiddleware(userRepo),
		AuthHandler:             handler.NewAuthHandler(authService),
		ApplicationHandler:      handler.NewApplicationHandler(applicationService),
		HeroSectionHandler:      handler.NewHeroSectionHandler(heroSectionService),
		FaqHandler:              handler.NewFaqHandler(faqService),
		GalleryHandler:          handler.NewGalleryHandler(galleryService),
		TestimonyHandler:        handler.NewTestimonyHandler(testimonyService),
		VillaPackageListHandler: handler.NewVillaPackageListHandler(villaPackageListService),
	}

	app.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	app.RegisterRoute()

	return app
}

func (a *App) Run(addr string) error {
	err := a.Router.Run(addr)
	return err
}

func setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			return name
		})
	}
}
