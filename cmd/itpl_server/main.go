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

	//init function 쓰기. 패키지 로드될 때 init function 실행됨.
	err = config.InitConfigs()

	if err != nil {
		panic("Failed to load configs: " + err.Error())
	}

	db, err := storage.InitMySQL(*config.DbCfg)
	if err != nil {
		panic("Failed to init mysql: " + err.Error())
	}

	storage.AutoMigrate(db)
	r := server.SetupRouter(db)

	r.Run(":8080")
}
