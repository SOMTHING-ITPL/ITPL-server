package course

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	UserId      uint    `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsAICreated bool    `json:"is_ai_created"`
}

type CourseDetail struct {
	gorm.Model
	CourseId uint `json:"course_id"`
	Day      int  `json:"day"`
	Sequence int  `json:"sequence"`
	PlaceId  uint `json:"place_id"`
}
