package main

import (
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	server "github.com/SOMTHING-ITPL/ITPL-server/internal/app"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	//set env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
