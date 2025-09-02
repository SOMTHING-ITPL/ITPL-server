package course

import (
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func CreateCourse(db *gorm.DB, user user.User, title string, description *string, facilityId uint) error {
	course := Course{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		IsAICreated: false, // by default
		FacilityID:  facilityId,
	}
	if err := db.Create(&course).Error; err != nil {
		return err
	}
	return nil
}

func GetCourseByCourseId(db *gorm.DB, courseId uint) (*Course, error) {
	var course Course
	err := db.Where("id = ?", courseId).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func GetCoursesByUserId(db *gorm.DB, userID uint) ([]Course, error) {
	var courses []Course
	err := db.Where("user_id = ?", userID).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func GetCourseDetails(db *gorm.DB, courseId uint) ([]CourseDetail, error) {
	var details []CourseDetail
	err := db.Where("course_id = ?", courseId).Order("day, sequence").Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func AddPlaceToCourse(db *gorm.DB, courseId uint, placeId uint, day int, sequence int) error {
	courseDetail := CourseDetail{
		CourseID: courseId,
		PlaceID:  placeId,
		Day:      day,
		Sequence: sequence,
	}
	if err := db.Create(&courseDetail).Error; err != nil {
		return err
	}
	return nil
}

func DeletePlaceFromCourse(db *gorm.DB, courseId uint, placeID uint) error {
	if err := db.Where("course_id = ? AND place_id = ?", courseId, placeID).Delete(&CourseDetail{}).Error; err != nil {
		return err
	}
	return nil
}

func ModifyCourse(db *gorm.DB, courseId uint, details []CourseDetail) error {
	if err := db.Where("Course_id = ?", courseId).Delete(&CourseDetail{}).Error; err != nil {
		return err
	}
	for _, detail := range details {
		detail.CourseID = courseId
		if err := db.Create(&detail).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetLastCoordinate(db *gorm.DB, course Course) place.Coordinate {
	details, err := GetCourseDetails(db, course.ID)
	if err != nil {
		defer log.Fatalf("failed to get course detail")
	}
	last := details[len(details)-1]
	lastPlace, err := place.GetPlaceById(db, last.ID)
	if err != nil {
		defer log.Fatalf("failed to get place")
	}
	return place.Coordinate{
		Latitude:  lastPlace.Latitude,
		Longitude: lastPlace.Longitude,
	}
}

func GetSpecificCouseDetail(db *gorm.DB, course Course, day, sequence int) CourseDetail {
	var detail CourseDetail
	err := db.Where("course_id = ? AND day = ? AND sequence = ?", course.ID, day, sequence).Find(&detail).Error
	if err != nil {
		defer log.Printf("failed to find specific course detail")
	}
	return detail
}
