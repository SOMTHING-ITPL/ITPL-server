package place

import (
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func buildReview(user user.User, text string, rating float64) review {
	return review{
		userId:  user.Username,
		rating:  rating,
		comment: text,
	}
}

func buildPlaceReview(placeId uint, rev review) PlaceReview {
	return PlaceReview{
		PlaceId: placeId,
		UserId:  rev.userId,
		Rating:  rev.rating,
		Comment: rev.comment,
	}
}

func addReviewToDB(db *gorm.DB, rev PlaceReview) {
	// insert review into the database
	if err := db.Create(&rev).Error; err != nil {
		panic("Failed to add review to database: " + err.Error())
	}
}

func WriteReview(db *gorm.DB, placeId uint, user user.User, text string, rating float64) error {
	addReviewToDB(db, buildPlaceReview(placeId, buildReview(user, text, rating)))
	return nil
}

func GetReviewInfo(db *gorm.DB, placeID uint) (ReviewInfo, error) {

	var result ReviewInfo

	err := db.Model(&PlaceReview{}).
		Select("COUNT(*) as count, IFNULL(AVG(rating), 0) as avg").
		Where("place_id = ?", placeID).
		Scan(&result).Error

	return result, err
}
