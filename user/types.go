package user

import (
	"time"

	"gorm.io/gorm"
)

type SocialProvider string

const (
	ProviderLocal  SocialProvider = "local"
	ProviderKakao  SocialProvider = "kakao"
	ProviderGoogle SocialProvider = "google"
)

type User struct {
	gorm.Model
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         string         `gorm:"type:varchar(100);not null" json:"user_id"`
	NickName       *string        `gorm:"type:varchar(100);default:null" json:"nickname,omitempty"`
	Email          *string        `gorm:"type:varchar(255);default:null" json:"email,omitempty"`
	SocialID       *string        `gorm:"type:varchar(255);index:idx_provider_social,unique;default:null" json:"social_id,omitempty"`
	SocialProvider SocialProvider `gorm:"type:enum('google','kakao','local');default:'local';not null;index:idx_provider_social,unique" json:"social_provider"`
	Photo          *string        `gorm:"type:varchar(255);default:null" json:"photo,omitempty"`
	EncryptPwd     *string        `gorm:"type:varchar(255);default:null" json:"encrypt_pwd,omitempty"`

	UserArtists []UserArtist `gorm:"foreignKey:UserID" json:"user_artists,omitempty"`
	UserGenres  []UserGenre  `gorm:"foreignKey:UserID" json:"user_genres,omitempty"`
}

type Artist struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	ImageURL  *string   `gorm:"type:varchar(255);default:null" json:"image_url,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	UserArtists []UserArtist `gorm:"foreignKey:ArtistID" json:"user_artists,omitempty"`
}

type Genre struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	UserGenres []UserGenre `gorm:"foreignKey:GenreID" json:"user_genres,omitempty"`
}

type UserArtist struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	ArtistID  uint      `gorm:"primaryKey" json:"artist_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Artist Artist `gorm:"foreignKey:ArtistID;constraint:OnDelete:CASCADE;" json:"artist"`
}

type UserGenre struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	GenreID   uint      `gorm:"primaryKey" json:"genre_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Genre Genre `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE;" json:"genre"`
}
