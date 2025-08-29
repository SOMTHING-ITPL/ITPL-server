package handler

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
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

// common res
type CommonRes struct {
	Message string
	Data    any
}

// performance common res
// 공연 목록
type PerformanceListRes struct {
	Performances []performanceShort `json:"performance"`
	Count        int                `json:"count"`
}

type FacilityListRes struct {
	Facilities []FacilityShort `json:"facility"`
	Count      int             `json:"count"`
}

// shortCut
type performanceShort struct {
	Id           uint   `json:"id"`
	Title        string `json:"title"`
	State        string `json:"state"`
	PosterURL    string `json:"poster_url"`
	FacilityName string `json:"facility_name"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

// detail
type PerformanceDetail struct {
	Id uint `json:"id"`

	Title     string    `json:"name"`       // prfnm
	StartDate time.Time `json:"start_date"` // prfpdfrom
	EndDate   time.Time `json:"end_date"`

	FacilityID   uint   `json:"facility_id"`
	FacilityName string `json:"facility_name"`

	AgeRating string `json:"age"`

	TicketPrice string `json:"price"`
	PosterURL   string `json:"poster"`

	Status       string  ` json:"state"`
	IsForeign    string  `json:"visit"` //내한 여부
	DateGuidance *string `json:"date_guidance"`

	IntroImageURL []performance.PerformanceImage      `json:"intro_url"`
	TicketSite    []performance.PerformanceTicketSite `json:"ticket_site"`
	LastModified  time.Time                           `json:"update_date"` // updatedate
}

// 공연 시설 목록
type FacilityShort struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	SeatCount string `json:"seat_count"` //없을 수도 있나?
}

type FacilityDetail struct {
	Id         uint
	Name       string  `json:"name"`
	OpenedYear *string `json:"open_year"`
	SeatCount  string  `json:"seat_count"`
	Phone      *string `json:"phone"`     // 전화번호
	Homepage   *string `json:"homepage" ` // 홈페이지
	Address    string  `json:"addree"`    // 주소
	Latitude   float64 `json:"latitude"`  // 위도
	Longitude  float64 `json:"longitude"` // 경도

	Restaurant string `json:"restaurant"`  // 음식점 유무
	Cafe       string `json:"cafe"`        // 카페 유무
	Store      string `json:"store"`       // 상점 유무
	ParkingLot string `json:"parking_lot"` // 주차시설
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

type ReviewImageResponse struct {
	URL string `json:"url"`
}

type PlaceReviewResponse struct {
	ID           uint                  `json:"id"`
	UserID       uint                  `json:"user_id"`
	UserNickname string                `json:"user_nickname"`
	Rating       float64               `json:"rating"`
	Comment      *string               `json:"comment"`
	Images       []ReviewImageResponse `json:"images"`
}
