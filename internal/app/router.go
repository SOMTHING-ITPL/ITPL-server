package server

import (
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/handler"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := user.NewRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	//this router does not needs auth
	public := r.Group("/")
	{
		//for health check
		registerHealthCheckRoutes(public)

		authGroup := public.Group("/auth")
		registerAuthRoutes(authGroup, userHandler)
	}

	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	// protected.Use(~)//should add middleWare
	{
		userGroup := protected.Group("/user")
		registerUserRoutes(userGroup, userHandler)

		courseGroup := protected.Group("/course")
		registerCourseRoutes(courseGroup)

		placeGroup := protected.Group("/place")
		registerPlaceRoutes(placeGroup, db, userRepo)

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
func registerAuthRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	rg.POST("/login", userHandler.LoginLocalUser())
	rg.POST("/check-id", userHandler.CheckValidId())

	rg.POST("/register", userHandler.RegisterLocalUser())
	rg.POST("/social-login", userHandler.LoginSocialUser())
}

// for about user
func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	rg.GET("/me", userHandler.GetUser)

	rg.GET("/artist", userHandler.GetArtists())
	rg.POST("/artist", userHandler.AddUserArtist())

	rg.GET("/genre", userHandler.GetGenres())
	rg.POST("/genre", userHandler.AddUserGenre())
}

// for about course
func registerCourseRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listCourseHandler)
}

// for about place
func registerPlaceRoutes(rg *gin.RouterGroup, db *gorm.DB, userRepo *user.Repository) {
	rg.GET("get-place-list", handler.GetPlaceList(db))
	rg.POST("/write-review", handler.WriteReviewHandler(db, userRepo))
	rg.POST("/delete-review/:rev_id", handler.DeleteReviewHandler(db, userRepo))
	// rg.POST("/", createPlaceHandler)
}

func registerConcertRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listConcertHandler)
	// rg.GET("/:id", getConcertHandler)
}
