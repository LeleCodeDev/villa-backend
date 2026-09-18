package bootstrap

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lelecodedev/villa-backend/internal/database"
	"github.com/lelecodedev/villa-backend/internal/handler"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/internal/service"
	"gorm.io/gorm"
)

type App struct {
	Router             *gin.Engine
	DB                 *gorm.DB
	HeroSectionHandler *handler.HeroSectionHandler
}

func NewApp() *App {
	r := gin.Default()
	db := database.NewDB()
	txManager := repository.NewTxManager(db)

	heroSectionRepo := repository.NewHeroSectionRepository(db)

	heroSectionService := service.NewHeroSectionService(txManager, heroSectionRepo)

	app := &App{
		Router:             r,
		DB:                 db,
		HeroSectionHandler: handler.NewHeroSectionHandler(heroSectionService),
	}

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
