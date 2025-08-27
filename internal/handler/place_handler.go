package handler

import (
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

		c.JSON(http.StatusOK, gin.H{"places": result})
	}
}
