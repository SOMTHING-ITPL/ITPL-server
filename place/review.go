package place

import (
	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func buildReview(user user.User, text string, rating float64, imgs []ReviewImage) review {
	return review{
		userId:     user.ID,
		nickname:   user.NickName,
		rating:     rating,
		comment:    &text,
		reviewImgs: imgs,
	}
}

func buildPlaceReview(placeId uint, rev review) PlaceReview {
	return PlaceReview{
		PlaceId:      placeId,
		UserId:       rev.userId,
		UserNickName: rev.nickname,
		Rating:       rev.rating,
		Comment:      rev.comment,
		Images:       rev.reviewImgs,
	}
}

func addReviewToDB(db *gorm.DB, rev PlaceReview) {
	// insert review into the database
	if err := db.Create(&rev).Error; err != nil {
		panic("Failed to add review to database: " + err.Error())
	}
}

func WriteReview(db *gorm.DB, placeId uint, user user.User, text string, rating float64, imgs []ReviewImage) error {
	addReviewToDB(db, buildPlaceReview(placeId, buildReview(user, text, rating, imgs)))
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

func GetPlaceReviews(db *gorm.DB, placeID uint) ([]PlaceReview, error) {
	var reviews []PlaceReview
	err := db.Preload("Images").Where("place_id = ?", placeID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func DeleteReview(db *gorm.DB, revId uint, bucketBasics aws.BucketBasics) error {
	// Load images
	var images []ReviewImage
	if err := db.Where("review_id = ?", revId).Find(&images).Error; err != nil {
		return err
	}

	// Delete images from S3
	for _, img := range images {
		err := aws.DeleteImage(bucketBasics.S3Client, bucketBasics.BucketName, img.Key)
		if err != nil {
			return err
		}
	}

	// 3. Delete review image records
	if err := db.Where("review_id = ?", revId).Delete(&ReviewImage{}).Error; err != nil {
		return err
	}

	// 4. Delete review record
	if err := db.Where("id = ?", revId).Delete(&PlaceReview{}).Error; err != nil {
		return err
	}

	return nil
}

func GetMyReviews(db *gorm.DB, userID uint) ([]PlaceReview, error) {
	var reviews []PlaceReview
	err := db.Preload("Images").Where("user_id = ?", userID).Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
