package performance

import (
	"time"

	"gorm.io/gorm"
)

// 공연 시설 테이블
type Facility struct {
	ID              uint64 `gorm:"primaryKey"`
	Name            string `gorm:"type:varchar(255);not null"`
	HallCount       *int
	Characteristics *string `gorm:"type:text"`
	OpenedYear      *int
	SeatCount       *int
	Phone           *string  `gorm:"type:varchar(50)"`
	Homepage        *string  `gorm:"type:varchar(255)"`
	Address         *string  `gorm:"type:varchar(255)"`
	Latitude        *float64 `gorm:"type:decimal(10,7)"`
	Longitude       *float64 `gorm:"type:decimal(10,7)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	Performances    []Performance  `gorm:"foreignKey:FacilityID"`
}

type Performance struct {
	ID            uint64 `gorm:"primaryKey"`
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type PerformanceTicketSite struct {
	PerformanceID uint64 `gorm:"primaryKey"`
	TicketSite    string `gorm:"primaryKey;type:varchar(255)"`
}
