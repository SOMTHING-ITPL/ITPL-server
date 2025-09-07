package artist

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

// artist 데이터 담는 부분 있어야 함.
type Artist struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null;unique" json:"name"`
	ImageKey string `gorm:"type:varchar(255)" json:"image_key,omitempty"`

	UserArtists []UserArtist `gorm:"foreignKey:ArtistID" json:"user_artists,omitempty"`
}

type UserArtist struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ArtistID  uint      `gorm:"primaryKey" json:"artist_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User   user.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Artist Artist    `gorm:"foreignKey:ArtistID;constraint:OnDelete:CASCADE;" json:"artist"`
}

type Repository struct {
	db *gorm.DB
}
