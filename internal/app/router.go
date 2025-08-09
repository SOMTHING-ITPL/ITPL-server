package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter sets up the router
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	//this router does not needs auth
	public := r.Group("/")
	{
		//for health check
		registerHealthCheckRoutes(public)

		authGroup := public.Group("/auth")
		registerAuthRoutes(authGroup)
	}

	protected := r.Group("/api")
	// protected.Use(~)//should add middleWare
	{
		userGroup := protected.Group("/user")
		registerUserRoutes(userGroup)

		courseGroup := protected.Group("/course")
		registerCourseRoutes(courseGroup)

		placeGroup := protected.Group("/place")
		registerPlaceRoutes(placeGroup)

		concertGroup := protected.Group("/concert")
		registerConcertRoutes(concertGroup)
	}

	return r
}

// This is for example for jo you can also check this for test
func registerHealthCheckRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

// for login & sign in
func registerAuthRoutes(rg *gin.RouterGroup) {
	// rg.GET("/:id", getUserHandler)
	// rg.PUT("/:id", updateUserHandler)
}

// for about user
func registerUserRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listPlaceHandler)
	// rg.POST("/", createPlaceHandler)
}

// for about course
func registerCourseRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listCourseHandler)
}

// for about place
func registerPlaceRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listPlaceHandler)
	// rg.POST("/", createPlaceHandler)
}

func registerConcertRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listConcertHandler)
	// rg.GET("/:id", getConcertHandler)
}
