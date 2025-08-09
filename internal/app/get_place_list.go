package server

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/gin-gonic/gin"
)

func getPlaceList(c *gin.Context) {
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

	places, err := place.LoadNearPlaces(coord, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"places": places})
}
