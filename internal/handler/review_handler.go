package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WriteReviewHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Review added successfully"})
	}
}

func DeleteReviewHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		revId, err := strconv.ParseUint(c.Param("id"), 10, 32)
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
