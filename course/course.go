package course

import (
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func CreateCourse(db *gorm.DB, user user.User, title string, description *string) error {
	course := Course{
		UserId:      user.ID,
		Title:       title,
		Description: description,
		IsAICreated: false, // by default
	}
	if err := db.Create(&course).Error; err != nil {
		return err
	}
	return nil
}

func GetCourseByCourseId(db *gorm.DB, courseID uint) (*Course, error) {
	var course Course
	err := db.Where("id = ?", courseID).First(&course).Error
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

func GetCourseDetails(db *gorm.DB, courseID uint) ([]CourseDetail, error) {
	var details []CourseDetail
	err := db.Where("course_id = ?", courseID).Order("day, sequence").Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func AddPlaceToCourse(db *gorm.DB, courseID uint, placeID uint, day int, sequence int) error {
	courseDetail := CourseDetail{
		CourseId: courseID,
		PlaceId:  placeID,
		Day:      day,
		Sequence: sequence,
	}
	if err := db.Create(&courseDetail).Error; err != nil {
		return err
	}
	return nil
}

func DeletePlaceFromCourse(db *gorm.DB, courseID uint, placeID uint) error {
	if err := db.Where("course_id = ? AND place_id = ?", courseID, placeID).Delete(&CourseDetail{}).Error; err != nil {
		return err
	}
	return nil
}
