package handler

import (
	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	aws_client "github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/calendar"
	"github.com/SOMTHING-ITPL/ITPL-server/email"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepository *user.Repository
	smtpRepository *email.Repository
	BucketBasics   *aws_client.BucketBasics
}

type PerformanceHandler struct {
	performanceRepo *performance.Repository
}

type PlaceHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	BucketBasics   *aws_client.BucketBasics
}

type CourseHandler struct {
	database        *gorm.DB
	userRepository  *user.Repository
	performanceRepo *performance.Repository
	bucketBasics    *aws_client.BucketBasics
}

type ChatRoomHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	bucketBasics   *aws_client.BucketBasics
}

type CalendarHandler struct {
	calendarRepo    *calendar.Repository
	performanceRepo *performance.Repository
}

type ArtistHandler struct {
	artistRepo   *artist.Repository
	BucketBasics *aws_client.BucketBasics
}
