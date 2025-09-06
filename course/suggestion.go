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

	restaurants, err := place.LoadNearPlaces(coord, 39, db, 3000)
	if err != nil {
		log.Printf("failed to load places")
	}
	random := rand.Intn(len(restaurants) - 1)

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

	restaurants, err := place.LoadNearPlaces(coord, 39, db, 3000)
	if err != nil {
		log.Printf("failed to load places")
	}
	random := rand.Intn(len(restaurants) - 1)

	//1일차 저녁
	restaurant := restaurants[random]
	AddPlaceToCourse(db, course.ID, restaurant.TourapiPlaceId, 1, 1)

	//숙소
	accommodations, err := place.LoadNearPlaces(coord, 32, db, 3000)
	if err != nil {
		log.Printf("failed to load places")
	}

	random = rand.Intn(len(accommodations) - 1)
	accommodation := accommodations[random]
	AddPlaceToCourse(db, course.ID, accommodation.TourapiPlaceId, 1, 2)

	acoord := place.Coordinate{
		Latitude:  accommodation.Latitude,
		Longitude: accommodation.Longitude,
	}

	//아점... ㅋㅋ
	restaurants, err = place.LoadNearPlaces(acoord, 39, db, 3000)
	random = rand.Intn(len(restaurants) - 1)
	AddPlaceToCourse(db, course.ID, restaurants[random].TourapiPlaceId, 2, 1)

	//관광지
	sights, err := place.LoadNearPlaces(acoord, 12, db, 3000)
	random = rand.Intn(len(sights) - 1)
	AddPlaceToCourse(db, course.ID, sights[random].TourapiPlaceId, 2, 2)

	//문화시설
	cultures, _ := place.LoadNearPlaces(coord, 14, db, 3000)
	random = rand.Intn(len(cultures) - 1)
  
	AddPlaceToCourse(db, course.ID, cultures[random].TourapiPlaceId, 2, 3)
	return course
}

func ThreeDayCourse(db *gorm.DB, user user.User, title string, description *string, facility performance.Facility) Course {

	course := TwoDayCourse(db, user, title, description, facility)
	coord := GetLastCoordinate(db, course)


	//2일차 저녁
	dinners, _ := place.LoadNearPlaces(coord, 39, db, 3000)
	random := rand.Intn(len(dinners) - 1)

	AddPlaceToCourse(db, course.ID, dinners[random].TourapiPlaceId, 2, 4)

	//2일차 숙소 - 1일차 숙소 그대로
	accommodation := GetSpecificCouseDetail(db, course, 1, 2)
	AddPlaceToCourse(db, course.ID, accommodation.PlaceID, 2, 5)

	coord = GetLastCoordinate(db, course)

	//3일차 아점
	brunchs, _ := place.LoadNearPlaces(coord, 39, db, 3000)
	AddPlaceToCourse(db, course.ID, brunchs[random].TourapiPlaceId, 3, 1)

	//마지막으로 쇼핑
	shoppings, _ := place.LoadNearPlaces(coord, 38, db, 3000)
  
	AddPlaceToCourse(db, course.ID, shoppings[random].TourapiPlaceId, 3, 2)
	return course
}
