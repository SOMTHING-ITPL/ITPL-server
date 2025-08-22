package handler

import (
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/course"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewCourseHandler(db *gorm.DB, userRepo *user.Repository) *CourseHandler {
	return &CourseHandler{
		database:       db,
		userRepository: userRepo,
	}
}

func (h *CourseHandler) CreateCourseHandler() func(c *gin.Context) {
	type req struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
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
		err = course.CreateCourse(h.database, user, request.Title, request.Description)
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
	return func(c *gin.Context) {
		courseId, err := strconv.ParseUint(c.Param("course_id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}

		details, err := course.GetCourseDetails(h.database, uint(courseId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get course details: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"course_details": details})
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
