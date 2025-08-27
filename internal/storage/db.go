package storage

import (
	"fmt"
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/SOMTHING-ITPL/ITPL-server/course"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// we will use gorm that is for make for handle database query easy
func InitMySQL(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MYSQL connection failed : %v", err)
		return nil, err
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&user.User{},
		&user.Artist{},
		&user.Genre{},
		&user.UserArtist{},
		&user.UserGenre{},
		&course.Course{},
		&course.CourseDetail{},
		&place.Place{},
		&place.ReviewImage{},
		&place.PlaceReview{},
		&performance.Facility{},
		&performance.Performance{},
		&performance.PerformanceTicketSite{},
		&performance.PerformanceImage{},

		// &performance.UserRecentPerformance{},
	)

	if err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}
}
