// Package bootstrap
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
		public.GET("/health", a.HealthHandler.CheckHealth)

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

		// Villa Package
		public.GET("/villa-packages", a.VillaPackageHandler.GetAllVillaPackages)

		// Villa Package List
		public.GET("/villa-package-lists", a.VillaPackageListHandler.GetAllVillaPackageLists)

		// Pricelist
		public.GET("/pricelists", a.PricelistHandler.GetAllPricelists)

		// Facility
		public.GET("/facilities", a.FacilityHandler.GetAllFacilities)

		// Specification
		public.GET("/specifications", a.SpecificationHandler.GetAllSpecifications)

		// Specification Gallery
		public.GET("/specification-gallery", a.SpecificationGalleryHandler.GetSpecificationGallery)

		// Addon
		public.GET("/addons", a.AddonHandler.GetAllAddons)

		// Addon List
		public.GET("/addon-lists", a.AddonListHandler.GetAllAddonLists)
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

		// Villa Package
		admin.GET("/villa-packages/:id", a.VillaPackageHandler.GetVillaPackageByID)
		admin.POST("/villa-packages", a.VillaPackageHandler.CreateVillaPackage)
		admin.PUT("/villa-packages/:id", a.VillaPackageHandler.UpdateVillaPackage)
		admin.DELETE("/villa-packages/:id", a.VillaPackageHandler.DeleteVillaPackage)
		admin.GET("/villa-packages/:id/lists", a.VillaPackageHandler.GetAllVillaPackageListsByID)

		// Villa Package List
		admin.GET("/villa-package-lists/:id", a.VillaPackageListHandler.GetVillaPackageListByID)
		admin.POST("/villa-package-lists", a.VillaPackageListHandler.CreateVillaPackageList)
		admin.PUT("/villa-package-lists/:id", a.VillaPackageListHandler.UpdateVillaPackageList)
		admin.DELETE("/villa-package-lists/:id", a.VillaPackageListHandler.DeleteVillaPackageList)

		// Pricelist
		admin.GET("/pricelists/:id", a.PricelistHandler.GetPricelistByID)
		admin.POST("/pricelists", a.PricelistHandler.CreatePricelist)
		admin.PUT("/pricelists/:id", a.PricelistHandler.UpdatePricelist)
		admin.DELETE("/pricelists/:id", a.PricelistHandler.DeletePricelist)

		// Facility
		admin.GET("/facilities/:id", a.FacilityHandler.GetFacilityByID)
		admin.POST("/facilities", a.FacilityHandler.CreateFacility)
		admin.PUT("/facilities/:id", a.FacilityHandler.UpdateFacility)
		admin.DELETE("/facilities/:id", a.FacilityHandler.DeleteFacility)

		// Specification
		admin.GET("/specifications/:id", a.SpecificationHandler.GetSpecificationByID)
		admin.POST("/specifications", a.SpecificationHandler.CreateSpecification)
		admin.PUT("/specifications/:id", a.SpecificationHandler.UpdateSpecification)
		admin.DELETE("/specifications/:id", a.SpecificationHandler.DeleteSpecification)

		// Specification Gallery
		admin.POST("/specification-gallery", a.SpecificationGalleryHandler.CreateSpecificationGallery)
		admin.PUT("/specification-gallery", a.SpecificationGalleryHandler.UpdateSpecificationGallery)
		admin.DELETE("/specification-gallery", a.SpecificationGalleryHandler.DeleteSpecificationGallery)

		// Addon
		admin.GET("/addons/:id", a.AddonHandler.GetAddonByID)
		admin.POST("/addons", a.AddonHandler.CreateAddon)
		admin.PUT("/addons/:id", a.AddonHandler.UpdateAddon)
		admin.DELETE("/addons/:id", a.AddonHandler.DeleteAddon)
		admin.GET("/addons/:id/lists", a.AddonHandler.GetAllAddonListsByID)

		// Addon List
		admin.GET("/addon-lists/:id", a.AddonListHandler.GetAddonListByID)
		admin.POST("/addon-lists", a.AddonListHandler.CreateAddonList)
		admin.PUT("/addon-lists/:id", a.AddonListHandler.UpdateAddonList)
		admin.DELETE("/addon-lists/:id", a.AddonListHandler.DeleteAddonList)
	}
}
