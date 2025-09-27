package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	aws_client "github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/gin-gonic/gin"
)

func NewArtistHandler(artistRepo *artist.Repository, bucket *aws_client.BucketBasics) *ArtistHandler {
	return &ArtistHandler{
		artistRepo:   artistRepo,
		BucketBasics: bucket,
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
			url, err := aws_client.GetPresignURL(h.BucketBasics.AwsConfig, h.BucketBasics.BucketName, a.ImageKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Image URL From AWS"})
				return
			}
			res = append(res, PreferSearchResponse{
				ID:       a.ID,
				Name:     a.Name,
				ImageUrl: url,
			})
		}
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(CommonRes{
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

func (h *ArtistHandler) PutArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file "})
			return
		}
		url, err := aws_client.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, fmt.Sprintf("artist/%s", name), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to upload profile image: %v", err),
			})
			return
		}

		artist := &artist.Artist{
			Name:     name,
			ImageKey: url,
		}

		if err = h.artistRepo.PutArtist(artist); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save artist "})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})

	}
}
