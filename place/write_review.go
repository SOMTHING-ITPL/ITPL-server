package place

import (
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func buildReview(user user.User, text string, rating float64) review {
	return review{
		userId:  user.NickName,
		rating:  rating,
		comment: text,
	}
}

func placeReview(placeId uint, rev review) PlaceReview {
	return PlaceReview{
		PlaceId: placeId,
		UserId:  rev.userId,
		Rating:  rev.rating,
		Comment: rev.comment,
	}
}

func addReviewToDB(db *gorm.DB, rev PlaceReview) {
	// insert review into the database
}

func WriteReview(db *gorm.DB, placeId uint, user user.User, text string, rating float64) error {
	addReviewToDB(db, placeReview(placeId, buildReview(user, text, rating)))
	return nil
}
