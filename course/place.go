package course

import (
	"log"
	"os"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/externalapi"

	"github.com/joho/godotenv"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Place struct {
	Tourapi_place_id string `json:"tourapi_place_id"`
	Category         int64  `json:"category"`
	Title            string `json:"title"`
	Address          string `json:"address"`
	Tel              string `json:"tel"`
	Longitude        string `json:"longitude"`
	Latitude         string `json:"latitude"`
	PlaceImage       string `json:"place_image"`
}

type Course struct {
	Places []Place `json:"course"`
}

func LoadNearPlaces(c Coordinate, category int64) ([]Place, error) {
	err := godotenv.Load("../.env")
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

	finalurl, err := externalapi.BuildURL(api_url, params)
	if err != nil {
		return nil, err
	}

	items, err := externalapi.FetchAndParseJSON(finalurl)

	if err != nil {
		return nil, err
	}

	var places []Place

	for _, item := range items {
		// Convert item to Place
		place := Place{
			Tourapi_place_id: item.ContentId,
			Category:         category,
			Title:            item.Title,
			Address:          item.Addr1 + " " + item.Addr2,
			Tel:              item.Tel,
			Longitude:        item.MapX,
			Latitude:         item.MapY,
			PlaceImage:       item.FirstImage,
		}
		places = append(places, place)
	}

	return places, nil
}
