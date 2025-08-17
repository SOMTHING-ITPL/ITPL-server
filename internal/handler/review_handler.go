package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WriteReviewHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Review added successfully"})
	}
}
