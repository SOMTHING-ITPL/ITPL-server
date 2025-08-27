package calendar

import (
	"os/user"

	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	PerformanceID uint `json:"performance_id"`
	UserID        uint `json:"user_id"`

	Year  uint `json:"year"`
	Month uint `json:"month"`
	Day   uint `json:"day"`

	Performance performance.Performance `gorm:"foreignKey:PerformanceID" json:"-"`
	User        user.User               `gorm:"foreignKey:UserID" json:"-"`
}
