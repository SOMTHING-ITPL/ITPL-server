package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/course"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCourseHandler(db *gorm.DB, userRepo *user.Repository, pRepo *performance.Repository) *CourseHandler {
	return &CourseHandler{
		database:        db,
		userRepository:  userRepo,
		performanceRepo: pRepo,
	}
}

func (h *CourseHandler) CreateCourseHandler() func(c *gin.Context) {
	type req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		FacilityID  uint    `json:"faciliry_id"`
	}
	return func(c *gin.Context) {
		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
		err = course.CreateCourse(h.database, user, request.Title, request.Description, request.FacilityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course created successfully"})
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
		c.JSON(http.StatusOK, gin.H{"message": "Place added to course successfully"})
	}
}

func (h *CourseHandler) GetCourseDetails() gin.HandlerFunc {
	type response struct {
		Course  course.Course
		Details []course.CourseDetail
	}

	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		co, err := course.GetCourseByCourseId(h.database, uint(courseId))
		details, err := course.GetCourseDetails(h.database, uint(courseId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course details: " + err.Error()})
			return
		}
		res := response{
			Course:  *co,
			Details: details,
		}

		c.JSON(http.StatusOK, gin.H{"course_details": res})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"courses": courses})
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

		c.JSON(http.StatusOK, gin.H{"message": "course modified successfully"})
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

		switch req.Days {
		case 1:
			createdCourse := course.OneDayCourse(h.database, me, "추천 코스", &desc, *facility)
			courseDetails, _ := course.GetCourseDetails(h.database, createdCourse.ID)
			res := response{
				Course:        createdCourse,
				CourseDetails: courseDetails,
			}
			c.JSON(http.StatusOK, gin.H{"createdCourse": res})
		default:
			c.JSON(http.StatusOK, gin.H{"message": "default"})
		}
	}
}
