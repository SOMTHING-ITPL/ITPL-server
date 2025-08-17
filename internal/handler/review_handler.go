package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WriteReviewHandler(db *gorm.DB, userRepo *user.Repository) gin.HandlerFunc {
	type req struct {
		PlaceId uint    `json:"place_id"`
		Text    string  `json:"text"`
		Rating  float64 `json:"rating"`
	}
	return func(c *gin.Context) {
		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		user, err := userRepo.GetById(userID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err := place.WriteReview(db, request.PlaceId, user, request.Text, request.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write review: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Review added successfully"})
	}
}

func DeleteReviewHandler(db *gorm.DB, userRepo *user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		revId, err := strconv.ParseUint(c.Param("rev_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}
		err = place.DeleteReview(db, uint(revId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}

func GetPlaceReviewsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		placeID, err := strconv.ParseUint(c.Param("place_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
			return
		}
		reviews, err := place.GetPlaceReviews(db, uint(placeID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}
