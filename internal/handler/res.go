package handler

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
)

type CommonRes struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// performance common res
// 공연 목록
type PerformanceListRes struct {
	Performances []performanceShort `json:"performance,omitempty"`
	Count        int                `json:"count"`
}

type FacilityListRes struct {
	Facilities []FacilityShort `json:"facility ,omitempty"`
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
	EndDate      string `json:"end_date "`
}

// detail
type PerformanceDetail struct {
	Id uint `json:"id"`

	Title     string    `json:"name"`       // prfnm
	StartDate time.Time `json:"start_date"` // prfpdfrom
	EndDate   time.Time `json:"end_date"`

	FacilityID   uint   `json:"facility_id"`
	FacilityName string `json:"facility_name"`

	AgeRating string `json:"age ,omitempty"`

	TicketPrice string `json:"price ,omitempty"`
	PosterURL   string `json:"poster ,omitempty"`

	Status       string  ` json:"state "`
	IsForeign    string  `json:"visit"` //내한 여부
	DateGuidance *string `json:"date_guidance"`

	IntroImageURL []string                            `json:"intro_url ,omitempty"`
	TicketSite    []performance.PerformanceTicketSite `json:"ticket_site ,omitempty"`
	LastModified  time.Time                           `json:"update_date ,omitempty"` // updatedate
}

// 공연 시설 목록
type FacilityShort struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	SeatCount string `json:"seat_count ,omitempty"` //없을 수도 있나?
}

type FacilityDetail struct {
	Id         uint    `json:"id"`
	Name       string  `json:"name"`
	OpenedYear *string `json:"open_year ,omitempty"`
	SeatCount  string  `json:"seat_count ,omitempty"`
	Phone      *string `json:"phone ,omitempty"`     // 전화번호
	Homepage   *string `json:"homepage ,omitempty" ` // 홈페이지
	Address    string  `json:"addree"`               // 주소
	Latitude   float64 `json:"latitude"`             // 위도
	Longitude  float64 `json:"longitude"`            // 경도

	Restaurant string `json:"restaurant ,omitempty"`  // 음식점 유무
	Cafe       string `json:"cafe ,omitempty"`        // 카페 유무
	Store      string `json:"store ,omitempty"`       // 상점 유무
	ParkingLot string `json:"parking_lot ,omitempty"` // 주차시설
}

type ReviewImageResponse struct {
	URL string `json:"url"`
}

type PlaceReviewResponse struct {
	ID           uint                  `json:"id"`
	PlaceID      uint                  `json:"place_id,omitempty"`
	PlaceName    string                `json:"place_name,omitempty"`
	UserID       uint                  `json:"user_id"`
	UserNickname string                `json:"user_nickname"`
	Rating       float64               `json:"rating"`
	Comment      *string               `json:"comment,omitempty"`
	CreatedAt    string                `json:"created_at"`
	Images       []ReviewImageResponse `json:"images,omitempty"`
}

type PlaceInfoResponse struct {
	PlaceInfo place.PlaceWithReview
	Reviews   []PlaceReviewResponse `json:"reviews,omitempty"`
}

type PreferSearchResponse struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}
