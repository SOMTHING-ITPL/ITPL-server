package handler

import (
	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/calendar"
	"github.com/SOMTHING-ITPL/ITPL-server/email"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepository *user.Repository
	smtpRepository *email.Repository
}

type PerformanceHandler struct {
	performanceRepo *performance.Repository
}

type PlaceHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	BucketBasics   *aws.BucketBasics
}

type CourseHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
}

type ChatRoomHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	smtpRepository *email.Repository
}

type CalendarHandler struct {
	calendarRepo    *calendar.Repository
	performanceRepo *performance.Repository
}
