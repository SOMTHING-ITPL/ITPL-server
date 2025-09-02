package course

import (
	"log"
	"math/rand"

	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func OneDayCourse(db *gorm.DB, user user.User, title string, description *string, facility performance.Facility) Course {
	course := Course{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		IsAICreated: true,
		FacilityID:  facility.ID,
	}
	if err := db.Create(&course).Error; err != nil {
		log.Printf("failed to create course")
	}

	coord := place.Coordinate{
		Latitude:  facility.Latitude,
		Longitude: facility.Longitude,
	}

	restaurants, err := place.LoadNearPlaces(coord, 39, db)
	if err != nil {
		log.Printf("failed to load places")
	}
	random := rand.Intn(10)

	restaurant := restaurants[random]
	AddPlaceToCourse(db, course.ID, restaurant.TourapiPlaceId, 1, 1)
	return course
}

func TwoDayCourse(db *gorm.DB, user user.User, title string, description *string, facility performance.Facility) Course {
	course := Course{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		IsAICreated: true,
		FacilityID:  facility.ID,
	}
	if err := db.Create(&course).Error; err != nil {
		log.Printf("failed to create course")
	}

	coord := place.Coordinate{
		Latitude:  facility.Latitude,
		Longitude: facility.Longitude,
	}

	restaurants, err := place.LoadNearPlaces(coord, 39, db)
	if err != nil {
		log.Printf("failed to load places")
	}
	random := rand.Intn(30)

	restaurant := restaurants[random]
	AddPlaceToCourse(db, course.ID, restaurant.TourapiPlaceId, 1, 1)
	return course
}
