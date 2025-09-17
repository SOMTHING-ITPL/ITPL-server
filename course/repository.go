package course

import (
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func CreateCourse(db *gorm.DB, user user.User, title string, description, imgKey *string, facilityId uint) error {
	course := Course{
		UserID:      user.ID,
		Title:       title,
		Description: description,
		IsAICreated: false, // by default
		FacilityID:  facilityId,
		ImageKey:    imgKey,
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

func GetLastCoordinate(db *gorm.DB, course Course) place.Coordinate {
	details, err := GetCourseDetails(db, course.ID)
	if err != nil {
		defer log.Printf("failed to get course detail")
	}
	last := details[len(details)]
	lastPlace, err := place.GetPlaceById(db, last.ID)
	if err != nil {
		defer log.Printf("failed to get place")
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

func AddPlaceToCourse(db *gorm.DB, courseId uint, placeId uint, day int, sequence int) error {
	place, err := place.GetPlaceById(db, placeId)
	if err != nil {
		defer log.Printf("place is not found")
	}
	courseDetail := CourseDetail{
		CourseID:   courseId,
		PlaceID:    placeId,
		Day:        day,
		Sequence:   sequence,
		PlaceTitle: place.Title,
		Address:    place.Address,
		Latitud:    place.Latitude,
		Longitude:  place.Longitude,
	}
	if err := db.Create(&courseDetail).Error; err != nil {
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

func ModifyCourseImageKey(db *gorm.DB, courseId uint, key *string) {
	course, err := GetCourseByCourseId(db, courseId)
	if err != nil {
		defer log.Printf("course does not exist : %v", err)
	}
	course.ImageKey = key
	db.Save(&course)
}

func DeletePlaceFromCourse(db *gorm.DB, courseId uint, placeID uint) error {
	if err := db.Where("course_id = ? AND place_id = ?", courseId, placeID).Delete(&CourseDetail{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCourse(db *gorm.DB, bucketBasics *aws.BucketBasics, courseId uint) error {

	deleteCourse, err := GetCourseByCourseId(db, courseId)
	if err != nil {
		return err
	}

	if err := db.Where("course_id = ?", courseId).Delete(&CourseDetail{}).Error; err != nil {
		return err
	}
	if err := db.Where("id = ?", courseId).Delete(&Course{}).Error; err != nil {
		return err
	}

	if err = aws.DeleteImage(bucketBasics.S3Client, bucketBasics.BucketName, *deleteCourse.ImageKey); err != nil {
		return err
	}

	return nil
}
