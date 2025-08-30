package handler

import (
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	"github.com/gin-gonic/gin"
)

func NewArtistHandler(artistRepo *artist.Repository) *ArtistHandler {
	return &ArtistHandler{
		artistRepo: artistRepo,
	}
}

func (h *ArtistHandler) GetArtists() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist, err := h.artistRepo.GetArtist()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get artist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": artist})
	}
}

func (h *ArtistHandler) AddUserArtist() gin.HandlerFunc {
	type req struct {
		ArtistIDs []uint `json:"artist_ids" binding:"required"`
	}
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "check request params"})
			return
		}

		if err := h.artistRepo.UpdateUserArtist(request.ArtistIDs, userID.(uint)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to set user artist on db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User artists updated successfully"})
	}
}

func (h *ArtistHandler) GetUserArtists() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		artists, err := h.artistRepo.GetUserArtists(userID.(uint))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Fail to get Artist"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": artists})

	}
}
