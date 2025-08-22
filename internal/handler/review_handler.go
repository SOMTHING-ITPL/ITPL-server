package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	database       *gorm.DB
	userRepository *user.Repository
}

func NewReviewHandler(db *gorm.DB, userRepo *user.Repository) *ReviewHandler {
	return &ReviewHandler{
		userRepository: userRepo,
		database:       db,
	}
}

func (h *ReviewHandler) WriteReviewHandler() gin.HandlerFunc {
	type req struct {
		PlaceId uint    `json:"place_id"`
		Text    string  `json:"text"`
		Rating  float64 `json:"rating"`
	}

	/*
		review image store logic here
	*/
	imgUrl := "test_url"

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
		user, err := h.userRepository.GetById(userID)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err := place.WriteReview(h.database, request.PlaceId, user, request.Text, request.Rating, imgUrl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write review: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Review added successfully"})
	}
}

func (h *ReviewHandler) DeleteReviewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		revId, err := strconv.ParseUint(c.Param("review_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)

		rev, err := place.GetReviewByID(h.database, uint(revId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
			return
		}

		if rev.UserId != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own reviews"})
			return
		}

		err = place.DeleteReview(h.database, uint(revId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}

func (h *ReviewHandler) GetPlaceReviewsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		placeID, err := strconv.ParseUint(c.Param("place_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
			return
		}
		reviews, err := place.GetPlaceReviews(h.database, uint(placeID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}

func (h *ReviewHandler) GetMyReviewsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)

		user, err := h.userRepository.GetById(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		reviews, err := place.GetMyReviews(h.database, user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}
