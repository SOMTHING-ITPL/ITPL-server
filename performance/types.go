package performance

import (
	"os/user"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// 공연 시설 테이블
type Facility struct {
	gorm.Model

	KopisFacilityKey string `gorm:"unique;not null" json:"kopis_facility_id"` // KOPIS 시설 키 (unique index)
	Name             string `gorm:"type:varchar(255);not null" json:"name"`   // 공연 시설명
	// HallCount        *int
	Characteristics *string `gorm:"type:text"` // 시설 특성
	OpenedYear      *string // 개관연도
	SeatCount       *string // 객석 수
	Phone           *string `gorm:"type:varchar(50)" json:"phone"`       // 전화번호
	Homepage        *string `gorm:"type:varchar(255)" json:"homepage" `  // 홈페이지
	Address         string  `gorm:"type:varchar(255)" json:"addree"`     // 주소
	Latitude        float64 `gorm:"type:decimal(10,7)" json:"latitude"`  // 위도
	Longitude       float64 `gorm:"type:decimal(10,7)" json:"longitude"` // 경도

	Restaurant string `gorm:"type:varchar(255)" json:"restaurant"`  // 음식점 유무
	Cafe       string `gorm:"type:varchar(255)" json:"cafe"`        // 카페 유무
	Store      string `gorm:"type:varchar(255)" json:"store"`       // 상점 유무
	ParkingLot string `gorm:"type:varchar(255)" json:"parking_lot"` // 주차시설

	Performances []Performance `gorm:"foreignKey:FacilityID"`
}

//LLM 에 포스터 던져주고 장르랑 keyword 공백으로 분리하여 저장해야 함.

type Performance struct {
	gorm.Model

	KopisPerformanceKey string    `gorm:"unique;not null" json:"kopis_performance_id"` // mt20id
	Title               string    `gorm:"type:varchar(255);not null" json:"name"`      // prfnm
	StartDate           time.Time `json:"start_date"`                                  // prfpdfrom
	EndDate             time.Time `json:"end_date"`
	// prfpdto
	KopisFacilityKey string `json:"kopis_facility_id"`

	FacilityID uint     `json:"facility_id"`
	Facility   Facility `gorm:"foreignKey:FacilityID" json:"facility"`

	Cast      *string `gorm:"type:text" json:"cast"`           // prfcast
	Crew      *string `gorm:"type:text" json:"crew"`           // prfcrew
	Runtime   *string `gorm:"type:varchar(50)" json:"runtime"` // prfruntime
	AgeRating *string `gorm:"type:varchar(50)" json:"age"`     // prfage
	// Producer      *string    `gorm:"type:varchar(255)" json:"producer"`  // entrpsnmP (제작사)
	// Organizer     *string    `gorm:"type:varchar(255)" json:"organizer"` // entrpsnmS (주관)
	// Sponsor       *string    `gorm:"type:varchar(255)" json:"sponsor"`   // entrpsnmH (주최)
	TicketPrice   *string   `gorm:"type:varchar(255)" json:"price"`    // pcseguidance
	PosterURL     *string   `gorm:"type:varchar(255)" json:"poster"`   // poster
	IntroImageURL *string   `gorm:"type:varchar(255)" json:"sty_urls"` // styurls>styurl
	Region        *string   `gorm:"type:varchar(100)" json:"area"`     // area
	Genre         *string   `gorm:"type:varchar(100)" json:"genre"`    // genrenm
	Status        *string   `gorm:"type:varchar(50)" json:"state"`     // prfstate
	IsForeign     string    `json:"visit"`                             // visit
	LastModified  time.Time `json:"update_date"`                       // updatedate
	Story         *string   `gorm:"type:text" json:"story"`            // sty
	DateGuidance  *string   `gorm:"type:text" json:"date_guidance"`    // dtguidance

	TicketSites []PerformanceTicketSite `gorm:"foreignKey:PerformanceID" json:"relates"` // 예매처
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
}

type PerformanceTicketSite struct {
	PerformanceID uint        `gorm:"primaryKey" json:"performance_id"`
	TicketSite    string      `gorm:"primaryKey;type:varchar(255)" json:"ticket_site"`
	Performance   Performance `gorm:"foreignKey:PerformanceID" json:"performance,omitempty"`
}

type UserRecentPerformance struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	UserID        uint64    `gorm:"not null;index:idx_user_viewed_at"`
	PerformanceID uint64    `gorm:"not null;uniqueIndex:ux_user_performance"`
	ViewedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_user_viewed_at"`

	User        user.User   `gorm:"foreignKey:UserID;references:ID"`
	Performance Performance `gorm:"foreignKey:PerformanceID;references:ID"`
}
