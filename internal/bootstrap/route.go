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
		public.GET("/hero-section", a.HeroSectionHandler.GetHeroSection)

		// FAQ
		public.GET("/faqs", a.FaqHandler.GetAllFaqs)

		// Gallery
		public.GET("/galleries", a.GalleryHandler.GetAllGalleries)

		// Application
		public.GET("/application", a.ApplicationHandler.GetApplication)
		// Testimony
		public.GET("/testimonies", a.TestimonyHandler.GetAllTestimonies)
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
		admin.POST("/galleries", a.GalleryHandler.CreateGallery)
		admin.PUT("/galleries/:id", a.GalleryHandler.UpdateGallery)
		admin.DELETE("/galleries/:id", a.GalleryHandler.DeleteGallery)

		// Application
		admin.POST("/application", a.ApplicationHandler.CreateApplication)
		admin.PUT("/application", a.ApplicationHandler.UpdateApplication)
		admin.DELETE("/application", a.ApplicationHandler.DeleteApplication)
		// Testimony
		admin.GET("/testimonies/:id", a.TestimonyHandler.GetTestimonyByID)
		admin.POST("/testimonies", a.TestimonyHandler.CreateTestimony)
		admin.PUT("/testimonies/:id", a.TestimonyHandler.UpdateTestimony)
		admin.DELETE("/testimonies/:id", a.TestimonyHandler.DeleteTestimony)
	}
}
