package place

import (
	"log"
	"os"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"

	"github.com/joho/godotenv"
)

func LoadNearPlaces(c Coordinate, category int64) ([]Place, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	api_url := os.Getenv("TOUR_API_URL") + "/locationBasedList2?"

	// for test, 수정 예정
	params := map[string]string{
		"serviceKey":    os.Getenv("SERVICE_KEY"),
		"numOfRows":     "30",
		"pageNo":        "1",
		"MobileOS":      "ETC",
		"MobileApp":     "AppTest",
		"_type":         "json",
		"arrange":       "E", // 거리순
		"mapX":          strconv.FormatFloat(c.Longitude, 'f', -1, 64),
		"mapY":          strconv.FormatFloat(c.Latitude, 'f', -1, 64),
		"radius":        "3000",
		"contentTypeId": strconv.FormatInt(category, 10),
		"lDongRegnCd":   "11",
		"lDongSignguCd": "140",
		"lclsSystm1":    "FD",
		"lclsSystm2":    "FD01",
		"lclsSystm3":    "FD010100",
		"areaCode":      "1",
		"sigunguCode":   "24",
		"cat1":          "A05",
		"cat2":          "A0502",
		"cat3":          "A05020100",
	}

	finalurl, err := api.BuildURL(api_url, params)
	if err != nil {
		return nil, err
	}

	items, err := api.FetchAndParseJSON(finalurl)

	if err != nil {
		return nil, err
	}

	var places []Place

	for _, item := range items {
		// Convert item to Place
		placeId, _ := strconv.ParseInt(item.ContentId, 10, 64)
		uPlaceId := uint(placeId)
		longitude, _ := strconv.ParseFloat(item.MapX, 64)
		latitude, _ := strconv.ParseFloat(item.MapY, 64)
		place := Place{
			TourapiPlaceId: uPlaceId,
			Category:       category,
			Title:          item.Title,
			Address:        item.Addr1 + " " + item.Addr2,
			Tel:            &item.Tel,
			Longitude:      longitude,
			Latitude:       latitude,
			PlaceImage:     &item.FirstImage,
		}
		places = append(places, place)
	}

	return places, nil
}
