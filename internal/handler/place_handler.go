package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPlaceHandler(db *gorm.DB, userRepo *user.Repository, bucketBasics *aws.BucketBasics) *PlaceHandler {
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

		category, err := strconv.ParseInt(c.Query("category"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid category": err})
			return
		}

		coord := place.Coordinate{
			Latitude:  lat,
			Longitude: lon,
		}

		places, err := place.LoadNearPlaces(coord, category, h.database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var result []place.PlaceWithReview
		for _, p := range places {
			reviewInfo, _ := place.GetReviewInfo(h.database, p.TourapiPlaceId)
			result = append(result, place.PlaceWithReview{
				Place:       p,
				ReviewCount: reviewInfo.Count,
				ReviewAvg:   reviewInfo.Avg,
			})
		}

		c.JSON(http.StatusOK, CommonRes{
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
				url, err := aws.GetPresignURL(
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
			})
		}

		placeInfoResponse := PlaceInfoResponse{
			PlaceInfo: pInfo,
			Reviews:   reviews,
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "Place Info",
			Data:    placeInfoResponse})
	}
}
