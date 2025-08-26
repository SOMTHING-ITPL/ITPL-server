package performance

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
}

// 공연 시설 테이블 자동으로 facility 를 가져오는ㄴ가 ? 테스트 해봐야 할 것 같은데 ..
type Facility struct {
	gorm.Model

	KopisFacilityKey string `gorm:"unique;not null"`                        // KOPIS 시설 키 (unique index)
	Name             string `gorm:"type:varchar(255);not null" json:"name"` // 공연 시설명
	// HallCount        *int
	Characteristics *string `gorm:"type:text"` // 시설 특성
	OpenedYear      *string `gorm:"type:varchar(255)" json:"open_year"`
	SeatCount       *string `gorm:"type:varchar(255)" json:"seat_count"`
	Phone           *string `gorm:"type:varchar(255)" json:"phone"`      // 전화번호
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

// LLM 에 포스터 던져주고 장르랑 keyword 공백으로 분리하여 저장해야 함.
// 장르는 ENUM ? 코드 1 2 3 4 5 ...로 구분해야 할 듯?
type Performance struct {
	gorm.Model

	KopisPerformanceKey string    `gorm:"unique;not null" `                       // mt20id
	Title               string    `gorm:"type:varchar(255);not null" json:"name"` // prfnm
	StartDate           time.Time `json:"start_date"`                             // prfpdfrom
	EndDate             time.Time `json:"end_date"`
	// prfpdto
	KopisFacilityKey string `json:"kopis_facility_id"`

	FacilityID   uint   `json:"facility_id"`
	FacilityName string `json:"facility_name"`

	Cast      *string `gorm:"type:text" json:"cast"`           // prfcast
	Crew      *string `gorm:"type:text" json:"crew"`           // prfcrew
	Runtime   *string `gorm:"type:varchar(50)" json:"runtime"` // prfruntime
	AgeRating *string `gorm:"type:varchar(50)" json:"age"`     // prfage
	// Producer      *string    `gorm:"type:varchar(255)" json:"producer"`  // entrpsnmP (제작사)
	// Organizer     *string    `gorm:"type:varchar(255)" json:"organizer"` // entrpsnmS (주관)
	// Sponsor       *string    `gorm:"type:varchar(255)" json:"sponsor"`   // entrpsnmH (주최)
	TicketPrice   *string `gorm:"type:varchar(255)" json:"price"`    // pcseguidance
	PosterURL     *string `gorm:"type:varchar(255)" json:"poster"`   // poster
	IntroImageURL *string `gorm:"type:varchar(255)" json:"sty_urls"` // styurls>styurl
	Region        *string `gorm:"type:varchar(100)" json:"area"`     // area
	// Genre         *string   `gorm:"type:varchar(100)" json:"genre"`    // genrenm
	Status       string    `gorm:"type:varchar(50)" json:"state"`  // prfstate
	IsForeign    string    `json:"visit"`                          // visit
	LastModified time.Time `json:"update_date"`                    // updatedate
	Story        *string   `gorm:"type:text" json:"story"`         // sty
	DateGuidance *string   `gorm:"type:text" json:"date_guidance"` // dtguidance
	Genre        int       `gorm:"type:int" json:"genre"`          // dtguidance
	Keyword      string    `gorm:"type:text" json:"keyword"`       // dtguidance

	TicketSites []PerformanceTicketSite `gorm:"foreignKey:PerformanceID"`
	Images      []PerformanceImage      `gorm:"foreignKey:PerformanceID"`

	// Facility Facility `gorm:"foreignKey:FacilityID" json:"facility"`
}

type PerformanceImage struct {
	PerformanceID uint        `gorm:"primaryKey" json:"performance_id"`
	URL           string      `gorm:"primaryKey;type:varchar(255)" json:"url"`
	Performance   Performance `gorm:"foreignKey:PerformanceID" json:"performance,omitempty"`
}

type PerformanceTicketSite struct {
	PerformanceID uint        `gorm:"primaryKey" json:"performance_id"`
	URL           string      `gorm:"primaryKey;type:varchar(255)" json:"url"`
	Name          string      `gorm:"type:varchar(255)" json:"name"`
	Performance   Performance `gorm:"foreignKey:PerformanceID" json:"performance,omitempty"`
}

type PerformanceUserLike struct {
	PerformanceID uint `gorm:"primaryKey" json:"performance_id"`
	UserID        uint `gorm:"primaryKey" json:"user_id"`

	Performance Performance `gorm:"foreignKey:PerformanceID" json:"performance,omitempty"`
	User        user.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// store in redis}
type PerformanceScore struct {
	ID    uint    `json:"id"`
	Score float64 `json:"score"`
}

// for preload
type PerformanceWithTicketsAndImage struct {
	Performance
	TicketSites       []PerformanceTicketSite `gorm:"foreignKey:PerformanceID"`
	PerformanceImages []PerformanceImage      `gorm:"foreignKey:PerformanceID"`
}
