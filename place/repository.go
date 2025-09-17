package place

import (
	"log"
	"os"
	"strconv"

	"gorm.io/gorm"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
)

func upsertByCreatedTime(db *gorm.DB, places []Place) error {
	for _, place := range places {
		var existingPlace Place
		result := db.Where("tourapi_place_id = ?", place.TourapiPlaceId).First(&existingPlace)

		switch result.Error {
		case nil:
			// Place exists, check if we need to update
			if existingPlace.CreatedTime != place.CreatedTime {
				updateResult := db.Model(&existingPlace).Updates(place)
				if updateResult.Error != nil {
					log.Printf("Failed to update place %d: %v", place.TourapiPlaceId, updateResult.Error)
				}
			}
			// no action needed if CreatedTime is the same
			break
		case gorm.ErrRecordNotFound:
			// Place does not exist, create a new one
			createResult := db.Create(&place)
			if createResult.Error != nil {
				log.Printf("Failed to create new place %d: %v", place.TourapiPlaceId, createResult.Error)
			}
			break
		default:
			// Some other error
			log.Printf("Database error while checking for place %d: %v", place.TourapiPlaceId, result.Error)
		}
	}
	return nil
}

func LoadNearPlaces(c Coordinate, category *string, db *gorm.DB, radius int64) ([]Place, error) {
	api_url := os.Getenv("TOUR_API_URL") + "/locationBasedList2?"
	params := map[string]string{
		"serviceKey": os.Getenv("SERVICE_KEY"),
		"numOfRows":  "30",
		"pageNo":     "1",
		"MobileOS":   "ETC",
		"MobileApp":  "AppTest",
		"_type":      "json",
		"arrange":    "E", // 거리순
		"mapX":       strconv.FormatFloat(c.Longitude, 'f', -1, 64),
		"mapY":       strconv.FormatFloat(c.Latitude, 'f', -1, 64),
		"radius":     strconv.FormatInt(radius, 10),
	}

	if category != nil && *category != "" { //category
		params["contentTypeId"] = *category
	}

	finalurl, err := api.BuildURL(api_url, params)
	if err != nil {
		return nil, err
	}
	log.Println("Final URL:", finalurl)

	items, err := api.FetchAndParseJSON(finalurl)

	if err != nil {
		return nil, err
	}

	var places []Place

	for _, item := range items {
		// Convert item to Place
		placeId, _ := strconv.ParseInt(item.ContentId, 10, 64)
		uPlaceId := uint(placeId)
		categoryID, _ := strconv.ParseInt(item.CategoryID, 10, 64)
		longitude, _ := strconv.ParseFloat(item.MapX, 64)
		latitude, _ := strconv.ParseFloat(item.MapY, 64)
		place := Place{
			TourapiPlaceId: uPlaceId,
			Category:       categoryID,
			Title:          item.Title,
			Address:        item.Addr1 + " " + item.Addr2,
			Tel:            &item.Tel,
			Longitude:      longitude,
			Latitude:       latitude,
			PlaceImage:     &item.FirstImage,
			CreatedTime:    item.CreatedTime,
		}
		places = append(places, place)
	}

	if err := upsertByCreatedTime(db, places); err != nil {
		log.Printf("Failed to upsert places: %v", err)
	}
	return places, nil
}

func GetPlaceById(db *gorm.DB, placeId uint) (*Place, error) {
	var place Place
	err := db.Where("tourapi_place_id = ?", placeId).First(&place).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Place not found
		}
		return nil, err // Other error
	}
	return &place, nil
}

func GetReviewInfo(db *gorm.DB, placeID uint) (ReviewInfo, error) {

	var result ReviewInfo

	err := db.Model(&PlaceReview{}).
		Select("COUNT(*) as count, IFNULL(AVG(rating), 0) as avg").
		Where("place_id = ?", placeID).
		Scan(&result).Error

	return result, err
}

func GetPlaceInfo(db *gorm.DB, placeID uint) (PlaceWithReview, error) {
	var reviews []PlaceReview
	err := db.Preload("Images").Where("place_id = ?", placeID).Find(&reviews).Error
	if err != nil {
		log.Printf("Failed to Load Review Images %v: ", err)
	}
	place, err := GetPlaceById(db, placeID)
	reviewInfo, err := GetReviewInfo(db, placeID)
	placeWithReview := PlaceWithReview{
		Place:       *place,
		ReviewCount: reviewInfo.Count,
		ReviewAvg:   reviewInfo.Avg,
	}

	return placeWithReview, err
}

func GetPlaceName(db *gorm.DB, placeID uint) (string, error) {
	var place Place
	err := db.Model(&Place{}).Where("tourapi_place_id = ?", placeID).Scan(&place).Error
	if err != nil {
		return "", err
	}

	return place.Title, nil
}

func GetImageByPlaceID(db *gorm.DB, placeID uint) string {
	place, err := GetPlaceById(db, placeID)
	if err != nil {
		defer log.Printf("Place does not exist")
	}
	return *place.PlaceImage
}

func SearchPlacesByTitle(db *gorm.DB, category int64, title string) ([]Place, error) {
	var places []Place
	if category == 0 {
		err := db.Where("title LIKE ?", "%"+title+"%").Find(&places).Error
		if err != nil {
			return nil, err
		}
	}

	err := db.Where("category = ? AND title LIKE ?", category, "%"+title+"%").Find(&places).Error
	if err != nil {
		return nil, err
	}
	return places, nil
}
