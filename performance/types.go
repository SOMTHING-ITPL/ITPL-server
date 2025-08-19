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

	KopisFacilityKey uint   `gorm:"unique;not null"`            // KOPIS 시설 키 (unique index)
	Name             string `gorm:"type:varchar(255);not null"` // 공연 시설명
	HallCount        *int
	Characteristics  *string `gorm:"type:text"` // 시설 특성
	OpenedYear       *int    // 개관연도
	SeatCount        *int    // 객석 수
	Phone            *string `gorm:"type:varchar(50)"`   // 전화번호
	Homepage         *string `gorm:"type:varchar(255)"`  // 홈페이지
	Address          string  `gorm:"type:varchar(255)"`  // 주소
	Latitude         float64 `gorm:"type:decimal(10,7)"` // 위도
	Longitude        float64 `gorm:"type:decimal(10,7)"` // 경도

	Performances []Performance `gorm:"foreignKey:FacilityID"`
}

type Performance struct {
	gorm.Model

	Title         string `gorm:"type:varchar(255);not null"`
	StartDate     *time.Time
	EndDate       *time.Time
	Cast          *string `gorm:"type:text"`
	Runtime       *string `gorm:"type:varchar(50)"`
	AgeRating     *string `gorm:"type:varchar(50)"`
	Producer      *string `gorm:"type:varchar(255)"`
	Organizer     *string `gorm:"type:varchar(255)"`
	Sponsor       *string `gorm:"type:varchar(255)"`
	TicketPrice   *string `gorm:"type:varchar(255)"`
	PosterURL     *string `gorm:"type:varchar(255)"`
	IntroImageURL *string `gorm:"type:varchar(255)"`
	Region        *string `gorm:"type:varchar(100)"`
	Genre         *string `gorm:"type:varchar(100)"`
	Status        *string `gorm:"type:varchar(50)"`
	IsForeign     *bool
	LastModified  *time.Time
	FacilityID    uint64
	Facility      Facility
	TicketSites   []PerformanceTicketSite `gorm:"foreignKey:PerformanceID"`
	DeletedAt     gorm.DeletedAt          `gorm:"index"`
}

type PerformanceTicketSite struct {
	PerformanceID uint64 `gorm:"primaryKey"`
	TicketSite    string `gorm:"primaryKey;type:varchar(255)"`
}

type UserRecentPerformance struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	UserID        uint64    `gorm:"not null;index:idx_user_viewed_at"`
	PerformanceID uint64    `gorm:"not null;uniqueIndex:ux_user_performance"`
	ViewedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_user_viewed_at"`

	User        user.User   `gorm:"foreignKey:UserID;references:ID"`
	Performance Performance `gorm:"foreignKey:PerformanceID;references:ID"`
}
