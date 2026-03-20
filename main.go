package main

import (
	"agnos-test/config"
	"agnos-test/database"
	"agnos-test/route"
)

func main() {
	cfg := config.Load()

	database.Migrate(cfg.DB)
	database.Seed(cfg.DB)

	r := route.SetupRoutes(cfg)
	r.Run(":" + cfg.ServerPort)
}
