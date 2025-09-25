package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	aws_client "github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPlaceHandler(db *gorm.DB, userRepo *user.Repository, bucketBasics *aws_client.BucketBasics) *PlaceHandler {
	return &PlaceHandler{
		database:       db,
		userRepository: userRepo,
		BucketBasics:   bucketBasics,
	}
}

func (h *PlaceHandler) GetPlaceList() gin.HandlerFunc {
	return func(c *gin.Context) {
		lat, err := strconv.ParseFloat(c.Query("latitude"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid latitude": err})
			return
		}

		lon, err := strconv.ParseFloat(c.Query("longitude"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid longitude": err})
			return
		}

		// category, err := strconv.ParseInt(c.Query("category"), 10, 64)
		category := c.Query("category")
		keyword := c.Query("keyword")

		coord := place.Coordinate{
			Latitude:  lat,
			Longitude: lon,
		}

		places, err := place.LoadNearPlaces(coord, category, h.database, 3000)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if keyword != "" {
			filtered := []place.Place{}
			trimmedKeyword := strings.TrimSpace(keyword)

			for _, p := range places {
				trimmedTitle := strings.TrimSpace(p.Title)

				if strings.Contains(trimmedTitle, trimmedKeyword) {
					filtered = append(filtered, p)
				}
			}
			places = filtered
		}

		result := make([]place.PlaceWithReview, 0)
		for _, p := range places {
			reviewInfo, _ := place.GetReviewInfo(h.database, p.TourapiPlaceId)
			result = append(result, place.PlaceWithReview{
				Place:       p,
				ReviewCount: reviewInfo.Count,
				ReviewAvg:   reviewInfo.Avg,
			})
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)

		_ = enc.Encode(CommonRes{
			Message: "Place Loaded Successfully",
			Data:    result,
		})
	}
}

func (h *PlaceHandler) GetPlaceInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		placeId, err := strconv.ParseUint(c.Param("place_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
			return
		}

		pInfo, err := place.GetPlaceInfo(h.database, uint(placeId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid place ID"})
			return
		}

		revs, err := place.GetPlaceReviews(h.database, uint(placeId))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err})
		}

		var reviews []PlaceReviewResponse
		for _, rev := range revs {
			var imgResponses []ReviewImageResponse
			for _, img := range rev.Images {
				url, err := aws_client.GetPresignURL(
					h.BucketBasics.AwsConfig,
					h.BucketBasics.BucketName,
					img.Key,
				)
				if err != nil {
					log.Printf("presign url error: %v", err)
					continue
				}
				imgResponses = append(imgResponses, ReviewImageResponse{URL: url})
			}

			reviews = append(reviews, PlaceReviewResponse{
				ID:           rev.ID,
				UserID:       rev.UserId,
				UserNickname: rev.UserNickName,
				Rating:       rev.Rating,
				Comment:      rev.Comment,
				Images:       imgResponses,
				CreatedAt:    rev.CreatedAt.Format(time.RFC3339),
			})
		}

		placeInfoResponse := PlaceInfoResponse{
			PlaceInfo: pInfo,
			Reviews:   reviews,
		}

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(c.Writer)
		enc.SetEscapeHTML(false)

		_ = enc.Encode(CommonRes{
			Message: "Place Info",
			Data:    placeInfoResponse,
		})
	}
}

func (h *PlaceHandler) SearchPlacesByTitleHandler() gin.HandlerFunc {
	type request struct {
		Title     string  `json:"title" binding:"required"`
		Category  int     `json:"category"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	return func(c *gin.Context) {

	}
}
