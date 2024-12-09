package main

import (
	"github.com/nulldiego/oh-back/api"
	"github.com/nulldiego/oh-back/config"
	"github.com/nulldiego/oh-back/internal/database"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	database.ConnectAutoMigrateDatabase()

	ohback := api.SetupApi()
	ohback.Run()
}
