package main

import (
	"github.com/SOMTHING-ITPL/ITPL-server/config"
	server "github.com/SOMTHING-ITPL/ITPL-server/internal/app"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
)

func main() {
	cfg, err := config.InitConfigs()

	if err != nil {
		panic("Failed to load configs: " + err.Error())
	}

	db, err := storage.InitMySQL(*cfg.DBConfig)
	if err != nil {
		panic("Failed to init mysql: " + err.Error())
	}

	storage.AutoMigrate(db)
	r := server.SetupRouter(db)

	r.Run(":8080")
}
