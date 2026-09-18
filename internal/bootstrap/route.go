package bootstrap

func (a *App) RegisterRoute() {
	api := a.Router.Group("/api")

	heroSection := api.Group("/hero-section")
	{
		heroSection.GET("", a.HeroSectionHandler.GetHeroSecton)
		heroSection.POST("", a.HeroSectionHandler.CreateHeroSection)
		heroSection.PUT("", a.HeroSectionHandler.UpdateHeroSection)
		heroSection.DELETE("", a.HeroSectionHandler.DeleteHeroSection)
	}

	faq := api.Group("/faqs")
	{
		faq.GET("", a.FaqHandler.GetAllFaqs)
		faq.GET("/:id", a.FaqHandler.GetFaqByID)
		faq.POST("", a.FaqHandler.CreateFaq)
		// faq.PUT("", a.HeroSectionHandler.UpdateHeroSection)
		// faq.DELETE("", a.HeroSectionHandler.DeleteHeroSection)
	}
}
