package main

import (
	server "github.com/SOMTHING-ITPL/ITPL-server/internal/app"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/storage"
)

func main() {
	storage.InitMySQL()
	storage.AutoMigrate()

	r := server.SetupRouter()

	r.Run(":8080")
}
