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

type Repository struct {
	db *gorm.DB
}

type User struct {
	gorm.Model

	//unique (조건부 unique 추가 email) 로컬 로그인용
	Email *string `gorm:"type:varchar(127);uniqueIndex:idx_local_email;default:null" json:"email,omitempty"`

	NickName string `gorm:"type:varchar(127);not null" json:"nickname"`

	//unique 소셜 로그인용
	SocialID       *string        `gorm:"type:varchar(255);default:null;uniqueIndex:idx_provider_social" json:"social_id,omitempty"`
	SocialProvider SocialProvider `gorm:"type:enum('google','kakao','local');default:'local';not null;uniqueIndex:idx_provider_social" json:"social_provider"`

	Photo      *string    `gorm:"type:varchar(255);default:null" json:"photo,omitempty"` //key 저장
	EncryptPwd *string    `gorm:"type:varchar(255);default:null" json:"encrypt_pwd,omitempty"`
	Birthday   *time.Time `gorm:"type:date;default:null" json:"birthday,omitempty"`

	//개인정보 약관
	Agree1 bool `gorm:"type:boolean;not null;default:false" json:"agree_personal_info,omitempty"`    // 필수 동의
	Agree2 bool `gorm:"type:boolean;not null;default:false" json:"agree_marketing,omitempty"`        // 마케팅 동의
	Agree3 bool `gorm:"type:boolean;not null;default:false" json:"agree_terms_of_service,omitempty"` // 서비스 약관 동의
	Agree4 bool `gorm:"type:boolean;not null;default:false" json:"agree_privacy_policy,omitempty"`   // 개인정보 처리방침 동의

	// UserArtists []UserArtist `gorm:"foreignKey:UserID" json:"user_artists,omitempty"`
	// UserGenres  []UserGenre  `gorm:"foreignKey:UserID" json:"user_genres,omitempty"`
}

type Genre struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null;unique" json:"name"`
	ImageKey string `gorm:"type:varchar(255)" json:"image_key,omitempty"`

	UserGenres []UserGenre `gorm:"foreignKey:GenreID" json:"user_genres,omitempty"`
}

type UserGenre struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	GenreID   uint      `gorm:"primaryKey" json:"genre_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Genre Genre `gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE;" json:"genre"`
}
