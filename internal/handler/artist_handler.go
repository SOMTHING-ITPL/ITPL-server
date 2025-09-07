package handler

import (
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	"github.com/SOMTHING-ITPL/ITPL-server/aws"
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
		//Get PresignedImage URL
		res := make([]PreferSearchResponse, 0, len(artist))
		for _, a := range artist {
			url, err := aws.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, a.ImageKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Image URL From AWS"})
				return
			}
			res = append(res, PreferSearchResponse{
				Name:     a.Name,
				ImageUrl: url,
			})
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    res,
		})
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

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
		})
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

		artistIDs := make([]uint, 0, len(artists))
		for _, g := range artists {
			artistIDs = append(artistIDs, g.ID)
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    artistIDs,
		})

	}
}
