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

	placeHandler := handler.NewPlaceHandler(db, userRepo)
	courseHandler := handler.NewCourseHandler(db, userRepo)
	chatRoomHandler := handler.NewChatRoomHandler(db, userRepo)

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
		registerCourseRoutes(courseGroup, courseHandler)

		placeGroup := protected.Group("/place")
		registerPlaceRoutes(placeGroup, placeHandler)

		concertGroup := protected.Group("/concert")
		registerConcertRoutes(concertGroup)

		chatGroup := protected.Group("/chat")
		registerChatRoutes(chatGroup, chatRoomHandler)
	}

	return r
}

// This is for example for jo you can also check this for test
func registerHealthCheckRoutes(rg *gin.RouterGroup) {
	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "I am very healthy"})
	})
}

// for login & sign in
func registerAuthRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	rg.POST("/login", userHandler.LoginLocalUser())
	rg.POST("/check-email", userHandler.CheckValidEmail())

	rg.POST("/register", userHandler.RegisterLocalUser())
	rg.POST("/social-login", userHandler.LoginSocialUser())
}

// for about user
func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	rg.GET("/me", userHandler.GetUser())
	rg.PATCH("/me", userHandler.UpdateProfile())

	rg.GET("/artist", userHandler.GetArtists())
	rg.POST("/artist", userHandler.AddUserArtist())
	rg.GET("/artist/me", userHandler.GetUserArtists())

	rg.GET("/genre", userHandler.GetGenres())
	rg.POST("/genre", userHandler.AddUserGenre())
	rg.GET("/genre/me", userHandler.GetUserGenres())
}

// for about course
func registerCourseRoutes(rg *gin.RouterGroup, courseHandler *handler.CourseHandler) {
	rg.POST("/create", courseHandler.CreateCourseHandler())
	rg.POST("/:course_id/place", courseHandler.AddPlaceToCourseHandler())
	rg.GET("/my-courses", courseHandler.GetMyCourses())
	rg.PATCH("/:course_id/details", courseHandler.ModifyCourseHandler())
}

// for about place
func registerPlaceRoutes(rg *gin.RouterGroup, placeHandler *handler.PlaceHandler) {
	rg.GET("/place-list", placeHandler.GetPlaceList())
	rg.GET("/reviews/:place_id", placeHandler.GetPlaceReviewsHandler())
	rg.POST("/review", placeHandler.WriteReviewHandler())
	rg.GET("/my-reviews", placeHandler.GetMyReviewsHandler())
	rg.DELETE("/review/:review_id", placeHandler.DeleteReviewHandler())
}

func registerConcertRoutes(rg *gin.RouterGroup) {
	// rg.GET("/", listConcertHandler)
	// rg.GET("/:id", getConcertHandler)
}

func registerChatRoutes(rg *gin.RouterGroup, chatRoomHandler *handler.ChatRoomHandler) {
	rg.POST("/room", chatRoomHandler.CreateChatRoom())
	// rg.GET("/chat-room/:id", getChatRoomHandler)
	// rg.POST("/chat-room/:id/message", postMessageHandler)
}
