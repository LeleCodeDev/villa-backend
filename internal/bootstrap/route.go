package bootstrap

func (a *App) RegisterRoute() {
	api := a.Router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", a.AuthHandler.Register)
		auth.POST("/login", a.AuthHandler.Login)
	}

	public := api.Group("")
	{
		public.GET("/hero-section", a.HeroSectionHandler.GetHeroSecton)

		public.GET("/faqs", a.FaqHandler.GetAllFaqs)
	}

	authenticated := api.Group("")
	authenticated.Use(a.AuthMiddleware)

	admin := authenticated.Group("")
	{
		admin.POST("/hero-section", a.HeroSectionHandler.CreateHeroSection)
		admin.PUT("/hero-section", a.HeroSectionHandler.UpdateHeroSection)
		admin.DELETE("/hero-section", a.HeroSectionHandler.DeleteHeroSection)

		admin.GET("/faqs/:id", a.FaqHandler.GetFaqByID)
		admin.POST("/faqs", a.FaqHandler.CreateFaq)
		admin.PUT("/faqs/:id", a.FaqHandler.UpdateFaq)
		admin.DELETE("/faqs/:id", a.FaqHandler.DeleteFaq)
	}
}
