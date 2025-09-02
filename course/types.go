package course

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	UserID      uint    `json:"user_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	IsAICreated bool    `json:"is_ai_created"`
	FacilityID  uint    `json:"facility_id"`
}

type CourseDetail struct {
	gorm.Model
	CourseID uint `json:"course_id"`
	Day      int  `json:"day"`
	Sequence int  `json:"sequence"`
	PlaceID  uint `json:"place_id"`
}
