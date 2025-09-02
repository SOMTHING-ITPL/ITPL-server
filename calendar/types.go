package calendar

import (
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	PerformanceID uint `json:"performance_id"`
	UserID        uint `json:"user_id"`

	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`

	Performance performance.Performance `gorm:"foreignKey:PerformanceID" json:"-"`
	User        user.User               `gorm:"foreignKey:UserID" json:"-"`
}

type Repository struct {
	db *gorm.DB
}
