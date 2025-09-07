package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/course"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCourseHandler(db *gorm.DB, userRepo *user.Repository, pRepo *performance.Repository, bucketBasics *aws.BucketBasics) *CourseHandler {
	return &CourseHandler{
		database:        db,
		userRepository:  userRepo,
		performanceRepo: pRepo,
		BucketBasics:    bucketBasics,
	}
}

func (h *CourseHandler) CreateCourseHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		description := c.PostForm("description")
		sfacilityId := c.PostForm("facility_id")
		facilityId, err := strconv.ParseUint(sfacilityId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid facitity id"})
			return
		}

		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		user, err := h.userRepository.GetById(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var imageKey *string
		file, err := c.FormFile("image")
		if err == nil {
			// 업로드 처리
			key, err := aws.UploadToS3(
				h.BucketBasics.S3Client,
				h.BucketBasics.BucketName,
				fmt.Sprintf("course_image/%d/%d", userID, facilityId),
				file,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
				return
			}
			imageKey = &key
		} else if !errors.Is(err, http.ErrMissingFile) {
			// 파일이 없는 경우는 무시, 그 외 에러만 처리
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
			return
		}
		err = course.CreateCourse(h.database, user, title, &description, imageKey, uint(facilityId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course: " + err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func (h *CourseHandler) AddPlaceToCourseHandler() gin.HandlerFunc {
	type req struct {
		PlaceId  uint `json:"place_id"`
		Day      int  `json:"day"`
		Sequence int  `json:"sequence"`
	}
	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		courseID := uint(courseId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if err := course.AddPlaceToCourse(h.database, courseID, request.PlaceId, request.Day, request.Sequence); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add place to course: " + err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func (h *CourseHandler) GetCourseDetails() gin.HandlerFunc {
	type response struct {
		Course  CourseInfoResponse
		Details []CourseDetailResponse
	}

	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		co, err := course.GetCourseByCourseId(h.database, uint(courseId))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "course does not exist"})
			return
		}

		courseInfo := CourseInfoResponse{
			ID:          co.ID,
			CreatedAt:   co.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   co.UpdatedAt.Format(time.RFC3339),
			UserID:      co.UserID,
			Title:       co.Title,
			Description: co.Description,
			IsAICreated: co.IsAICreated,
			FacilityID:  co.FacilityID,
			ImageKey:    co.ImageKey,
		}
		details, err := course.GetCourseDetails(h.database, uint(courseId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course details: "})
			return
		}
		var courseDetailInfos []CourseDetailResponse
		for _, detail := range details {
			courseDetailInfos = append(courseDetailInfos, CourseDetailResponse{
				ID:         detail.ID,
				CreatedAt:  detail.CreatedAt.Format(time.RFC3339),
				UpdatedAt:  detail.UpdatedAt.Format(time.RFC3339),
				CourseID:   detail.CourseID,
				Day:        detail.Day,
				Sequence:   detail.Sequence,
				PlaceID:    detail.PlaceID,
				PlaceTitle: detail.PlaceTitle,
				Address:    detail.Address,
				Latitud:    detail.Latitud,
				Longitude:  detail.Longitude,
			})
		}

		res := response{
			Course:  courseInfo,
			Details: courseDetailInfos,
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "Course Details",
			Data:    res,
		})
	}
}

func (h *CourseHandler) GetMyCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)

		courses, err := course.GetCoursesByUserId(h.database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses: "})
			return
		}
		var courseInfos []CourseInfoResponse
		for _, course := range courses {
			courseInfos = append(courseInfos, CourseInfoResponse{
				ID:          course.ID,
				CreatedAt:   course.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   course.UpdatedAt.Format(time.RFC3339),
				UserID:      userID,
				Title:       course.Title,
				Description: course.Description,
				IsAICreated: course.IsAICreated,
				FacilityID:  course.FacilityID,
				ImageKey:    course.ImageKey,
			})
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "My Courses",
			Data:    courseInfos,
		})
	}
}

func (h *CourseHandler) ModifyCourseHandler() gin.HandlerFunc {
	type request struct {
		Details []course.CourseDetail `json:"details" binding:"required"`
	}
	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		var req request

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := course.ModifyCourse(h.database, uint(courseId), req.Details); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (h *CourseHandler) CourseSuggestionHandler() gin.HandlerFunc {
	type request struct {
		FacilityID uint `json:"facility_id"`
		Days       uint `json:"days"`
	}
	type response struct {
		Course        course.Course
		CourseDetails []course.CourseDetail
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		me, err := h.userRepository.GetById(userID)
		facility, err := h.performanceRepo.GetFacilityById(req.FacilityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		desc := fmt.Sprintf("%s 님을 위한 코스입니다.", me.NickName)

		var createdCourse course.Course
		switch req.Days {
		case 1:
			createdCourse = course.OneDayCourse(h.database, me, "추천 코스", &desc, *facility)
			break
		case 2:

			createdCourse = course.TwoDayCourse(h.database, me, "추천 코스", &desc, *facility)
			break
		case 3:
			createdCourse = course.ThreeDayCourse(h.database, me, "추천 코스", &desc, *facility)
			break

		default:
			c.JSON(http.StatusOK, gin.H{"message": "cannot generate course"})
			return
		}

		courseDetails, _ := course.GetCourseDetails(h.database, createdCourse.ID)
		res := response{
			Course:        createdCourse,
			CourseDetails: courseDetails,
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "Course Created",
			Data:    res,
		})
	}
}

func (h *CourseHandler) ModifyCourseImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get course ID
		courseID := c.PostForm("course_id")
		icourseID, err := strconv.ParseUint(courseID, 10, 32)
		ucourseID := uint(icourseID)

		// Upload image to S3
		file, err := c.FormFile("image") // images -> image (단수)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
			return
		}
		key, err := aws.UploadToS3(h.BucketBasics.S3Client, h.BucketBasics.BucketName, fmt.Sprintf("course_images/%d", ucourseID), file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}

		// Set image key
		course.ModifyCourseImageKey(h.database, ucourseID, &key)

		c.Status(http.StatusNoContent)
	}
}

func (h *CourseHandler) DeleteCourseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		err = course.DeleteCourse(h.database, uint(courseId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete course"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
