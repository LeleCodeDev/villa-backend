package bootstrap

import (
	"github.com/lelecodedev/villa-backend/internal/middleware"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func (a *App) RegisterRoute() {
	api := a.Router.Group("/api")

	api.Static("/uploads", "uploads")

	auth := api.Group("/auth")
	{
		auth.POST("/register", a.AuthHandler.Register)
		auth.POST("/login", a.AuthHandler.Login)
	}

	public := api.Group("")
	{
		// Hero Section
		public.GET("/hero-section", a.HeroSectionHandler.GetHeroSecton)

		// FAQ
		public.GET("/faqs", a.FaqHandler.GetAllFaqs)

		// Gallery
		public.GET("/galleries", a.GalleryHandler.GetAllGalleries)
	}

	authenticated := api.Group("")
	authenticated.Use(a.AuthMiddleware)

	admin := authenticated.Group("")
	admin.Use(middleware.RoleMiddleware(model.RoleAdmin))
	{
		// Hero Section
		admin.POST("/hero-section", a.HeroSectionHandler.CreateHeroSection)
		admin.PUT("/hero-section", a.HeroSectionHandler.UpdateHeroSection)
		admin.DELETE("/hero-section", a.HeroSectionHandler.DeleteHeroSection)

		// FAQ
		admin.GET("/faqs/:id", a.FaqHandler.GetFaqByID)
		admin.POST("/faqs", a.FaqHandler.CreateFaq)
		admin.PUT("/faqs/:id", a.FaqHandler.UpdateFaq)
		admin.DELETE("/faqs/:id", a.FaqHandler.DeleteFaq)

		// Gallery
		admin.GET("/galleries/:id", a.GalleryHandler.GetGalleryByID)
		admin.POST("/galleries/", a.GalleryHandler.CreateGallery)
	}
}
