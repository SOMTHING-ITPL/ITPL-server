package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/course"
	_ "github.com/SOMTHING-ITPL/ITPL-server/internal/externalapi"
)

type placeListRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  int64   `json:"category"`
}

func GetPlaceList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
		lon, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
		category, err := strconv.ParseInt(r.URL.Query().Get("category"), 10, 64)

		coord := course.Coordinate{
			Latitude:  lat,
			Longitude: lon,
		}

		places, err := course.LoadNearPlaces(coord, category)
		if err != nil {
			http.Error(w, "failed to load places: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"places": places,
		})
	}
}
