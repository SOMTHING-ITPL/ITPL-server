package storage

import (
	"fmt"
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//we will use gorm that is for make for handle database query easy

var DB *gorm.DB

func InitMySQL() {
	cfg := config.GetDBConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MYSQL connection failed : %v", err)
	}

	DB = db
}

func AutoMigrate() {
	err := DB.AutoMigrate(
	// &model.User{}, have to add model struct
	)

	if err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}
}
