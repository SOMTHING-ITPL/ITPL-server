package user

import "gorm.io/gorm"

type SocialProvider string

const (
	ProviderLocal  SocialProvider = "local"
	ProviderKakao  SocialProvider = "kakao"
	ProviderGoogle SocialProvider = "google"
)

// User represents a registered user in the system.
// Social login is supported via Google and Kakao.
// - Local users have EncryptPwd set and SocialProvider = "local"
// - Social users have SocialProvider != "local" and EncryptPwd = NULL
// gorm 주석이 원래 저럼? .... 으....
type User struct {
	gorm.Model
	Username       string         `gorm:"type:varchar(100);not null" json:"username"`
	NickName       *string        `gorm:"type:varchar(100);default:null" json:"nickname,omitempty"`
	Email          *string        `gorm:"type:varchar(255);default:null" json:"email,omitempty"`
	SocialID       *string        `gorm:"type:varchar(255);index:idx_provider_social,unique;default:null" json:"social_id,omitempty"`
	SocialProvider SocialProvider `gorm:"type:enum('google','kakao','local');default:'local';not null;index:idx_provider_social,unique" json:"social_provider"`
	Photo          *string        `gorm:"type:varchar(255);default:null" json:"photo,omitempty"`
	EncryptPwd     *string        `gorm:"type:varchar(255);default:null" json:"encrypt_pwd,omitempty"`
}
