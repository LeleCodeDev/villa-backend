package bootstrap

func (a *App) RegisterRoute() {
	api := a.Router.Group("/api")

	heroSection := api.Group("/hero-section")
	{
		heroSection.GET("", a.HeroSectionHandler.GetHeroSecton)
		heroSection.POST("", a.HeroSectionHandler.CreateHeroSection)
		heroSection.PUT("", a.HeroSectionHandler.UpdateHeroSection)
	}
}
