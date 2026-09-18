package main

import (
	"fmt"

	"github.com/lelecodedev/villa-backend/internal/bootstrap"
	"github.com/lelecodedev/villa-backend/internal/config"
)

func init() {
	config.LoadConfig()
}

func main() {
	app := bootstrap.NewApp()

	if err := app.Run(":" + config.Env.Port); err != nil {
		panic(fmt.Sprintf("Application failed to run %v", err.Error()))
	}
}
