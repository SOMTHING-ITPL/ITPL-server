package handler

import (
	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/calendar"
	"github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/SOMTHING-ITPL/ITPL-server/email"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepository *user.Repository
	smtpRepository *email.Repository
	BucketBasics   *s3.BucketBasics
}

type PerformanceHandler struct {
	performanceRepo *performance.Repository
}

type PlaceHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	BucketBasics   *s3.BucketBasics
}

type CourseHandler struct {
	database        *gorm.DB
	userRepository  *user.Repository
	performanceRepo *performance.Repository
	bucketBasics    *s3.BucketBasics
}

type ChatRoomHandler struct {
	chatRoomRepository *chat.ChatRoomRepository
	userRepository     *user.Repository
	bucketBasics       *s3.BucketBasics
	tableBasics        *dynamo.TableBasics
}

type CalendarHandler struct {
	calendarRepo    *calendar.Repository
	performanceRepo *performance.Repository
}

type ArtistHandler struct {
	artistRepo   *artist.Repository
	BucketBasics *s3.BucketBasics
}

type ChatHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
	bucketBasics   *s3.BucketBasics
	tableBasics    *dynamo.TableBasics
}
