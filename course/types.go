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
	ImageKey    *string `json:"image_key,omitempty"`
}

type CourseDetail struct {
	gorm.Model

	CourseID   uint    `json:"course_id"`
	Day        int     `json:"day"`
	Sequence   int     `json:"sequence"`
	PlaceID    uint    `json:"place_id"`
	PlaceTitle string  `json:"place_title"`
	Address    string  `json:"address"`
	Latitud    float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}
