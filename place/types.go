package place

import (
	"gorm.io/gorm"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Place struct {
	TourapiPlaceId uint    `json:"tourapi_place_id" gorm:"column:tourapi_place_id" gorm:"primaryKey"`
	Category       int64   `json:"category" gorm:"column:category"`
	Title          string  `json:"title" gorm:"column:title"`
	Address        string  `json:"address" gorm:"column:address"`
	Tel            *string `json:"tel" gorm:"column:tel"`
	Longitude      float64 `json:"longitude" gorm:"column:longitude"`
	Latitude       float64 `json:"latitude" gorm:"column:latitude"`
	PlaceImage     *string `json:"place_image" gorm:"column:place_image"`
	CreatedTime    string  `json:"createdtime" gorm:"column:created_time"`
}

type ReviewInfo struct {
	Count int64
	Avg   float64
}

type PlaceWithReview struct {
	Place
	ReviewCount int64   `json:"review_count"`
	ReviewAvg   float64 `json:"review_avg"`
}

type PlaceReview struct {
	gorm.Model
	PlaceId      uint    `json:"place_id" gorm:"column:place_id"`
	UserId       uint    `json:"user_id" gorm:"column:user_id"`
	UserNickName string  `json:"user_nickname" gorm:"column:user_nickname"`
	Rating       float64 `json:"rating" gorm:"column:rating"`
	Comment      *string `json:"comment" gorm:"column:comment"`
	ReviewImage  *string `json:"review_image" gorm:"column:review_image"`
}

type review struct {
	userId    uint
	nickname  string
	rating    float64
	comment   *string
	reviewImg *string
}
