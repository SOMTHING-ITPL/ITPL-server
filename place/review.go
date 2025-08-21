package place

import (
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func buildReview(user user.User, text string, rating float64, img string) review {
	return review{
		userId:    user.ID,
		nickname:  user.NickName,
		rating:    rating,
		comment:   &text,
		reviewImg: &img,
	}
}

func buildPlaceReview(placeId uint, rev review) PlaceReview {
	return PlaceReview{
		PlaceId:      placeId,
		UserId:       rev.userId,
		UserNickName: rev.nickname,
		Rating:       rev.rating,
		Comment:      rev.comment,
		ReviewImage:  rev.reviewImg,
	}
}

func addReviewToDB(db *gorm.DB, rev PlaceReview) {
	// insert review into the database
	if err := db.Create(&rev).Error; err != nil {
		panic("Failed to add review to database: " + err.Error())
	}
}

func WriteReview(db *gorm.DB, placeId uint, user user.User, text string, rating float64, img string) error {
	addReviewToDB(db, buildPlaceReview(placeId, buildReview(user, text, rating, img)))
	return nil
}

func GetReviewByID(db *gorm.DB, revId uint) (*PlaceReview, error) {
	var review PlaceReview
	err := db.Where("id = ?", revId).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func GetReviewInfo(db *gorm.DB, placeID uint) (ReviewInfo, error) {

	var result ReviewInfo

	err := db.Model(&PlaceReview{}).
		Select("COUNT(*) as count, IFNULL(AVG(rating), 0) as avg").
		Where("place_id = ?", placeID).
		Scan(&result).Error

	return result, err
}

func GetPlaceReviews(db *gorm.DB, placeID uint) ([]PlaceReview, error) {
	var reviews []PlaceReview
	err := db.Where("place_id = ?", placeID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func DeleteReview(db *gorm.DB, revId uint) error {
	err := db.Where("id = ?", revId).Delete(&PlaceReview{}).Error
	if err != nil {
		return err
	}
	return nil
}

func GetMyReviews(db *gorm.DB, userID uint) ([]PlaceReview, error) {
	var reviews []PlaceReview
	err := db.Where("user_id = ?", userID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
