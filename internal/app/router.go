package server

import (
	"fmt"
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/artist"
	"github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/SOMTHING-ITPL/ITPL-server/calendar"
	"github.com/SOMTHING-ITPL/ITPL-server/email"
	"github.com/SOMTHING-ITPL/ITPL-server/internal/handler"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redisDB *redis.Client, bucketBasics *aws.BucketBasics) *gin.Engine {
	logger, err := NewLogger()
	if err != nil {
		fmt.Printf("fail to get logger %s", err)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(logger))

	//이거 좀 너무 더러운데 어케 정리하고 싶음.
	userRepo := user.NewRepository(db)
	smtpRepo := email.NewRepository(redisDB)
	performanceRepo := performance.NewRepository(db, redisDB)
	calendarRepo := calendar.NewRepository(db)
	artistRepo := artist.NewRepository(db)

	userHandler := handler.NewUserHandler(userRepo, smtpRepo, bucketBasics)
	performanceHandler := handler.NewPerformanceHandler(performanceRepo)
	courseHandler := handler.NewCourseHandler(db, userRepo, performanceRepo, bucketBasics)
	placeHandler := handler.NewPlaceHandler(db, userRepo, bucketBasics)
	calendarHandler := handler.NewCalendarHandler(calendarRepo, performanceRepo)
	artistHandler := handler.NewArtistHandler(artistRepo, bucketBasics)
	chatRoomHandler := handler.NewChatRoomHandler(db, userRepo, bucketBasics)

	//this router does not needs auth
	public := r.Group("/api")
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
		// putData := protected.Group("/testdata/sample")
		// putData.POST("/artist", artistHandler.PutArtist())
		// putData.POST("genre", userHandler.PutGenre())

		userGroup := protected.Group("/user")
		registerUserRoutes(userGroup, userHandler)
		registerArtistRoutes(userGroup, artistHandler)
		registerCalendarRoutes(userGroup, calendarHandler)
		registerUserPerformanceRoutes(userGroup, performanceHandler)

		courseGroup := protected.Group("/course")
		registerCourseRoutes(courseGroup, courseHandler)

		placeGroup := protected.Group("/place")
		registerPlaceRoutes(placeGroup, placeHandler)

		performanceGroup := protected.Group("/performance")
		registerPerformanceRoutes(performanceGroup, performanceHandler)

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
	rg.POST("/check-email", userHandler.SendEmailCode())
	rg.POST("/verify-email", userHandler.VerifyEmailCode())
	rg.POST("/register-local", userHandler.RegisterLocalUser())
	rg.POST("/social-login", userHandler.LoginSocialUser())
}

func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	rg.GET("/me", userHandler.GetUser())
	rg.PATCH("/me", userHandler.UpdateProfile())
	rg.GET("/genre", userHandler.GetGenres())
	rg.POST("/genre", userHandler.AddUserGenre())
	rg.GET("/genre/me", userHandler.GetUserGenres())
}
func registerArtistRoutes(rg *gin.RouterGroup, artistHandler *handler.ArtistHandler) {

	rg.GET("/artist", artistHandler.GetArtists())
	rg.POST("/artist", artistHandler.AddUserArtist())
	rg.GET("/artist/me", artistHandler.GetUserArtists())

}
func registerCalendarRoutes(rg *gin.RouterGroup, calendarHandler *handler.CalendarHandler) {
	rg.GET("/calendar", calendarHandler.GetCalendarData())
	rg.POST("/calendar", calendarHandler.CreateCalendarData())
	rg.DELETE("/calendar/:id", calendarHandler.DeleteCalendarData())
}

func registerUserPerformanceRoutes(rg *gin.RouterGroup, performanceHandler *handler.PerformanceHandler) {
	rg.GET("/performance/", performanceHandler.GetPerformanceLike()) //유저 Performance 조회
	rg.POST("/performance/:id", performanceHandler.CreatePerformanceLike())
	rg.DELETE("/performance/:id", performanceHandler.DeletePerformanceLike())

	rg.GET("/performance/recent", performanceHandler.GetRecentViewPerformance())

}

// for about course
func registerCourseRoutes(rg *gin.RouterGroup, courseHandler *handler.CourseHandler) {
	rg.POST("/create", courseHandler.CreateCourseHandler())

	rg.GET("/my-courses", courseHandler.GetMyCourses())
	rg.GET(":course_id/details", courseHandler.GetCourseDetails())

	rg.POST("/:course_id/place", courseHandler.AddPlaceToCourseHandler())
	rg.PATCH("/:course_id/details", courseHandler.ModifyCourseHandler())
	rg.PATCH("/:course_id/image", courseHandler.ModifyCourseImage())
	rg.POST("/suggestion", courseHandler.CourseSuggestionHandler())

	rg.DELETE("/:course_id", courseHandler.DeleteCourseHandler())
}

// for about place
func registerPlaceRoutes(rg *gin.RouterGroup, placeHandler *handler.PlaceHandler) {
	rg.GET("/list", placeHandler.GetPlaceList())
	rg.GET("/:place_id/info", placeHandler.GetPlaceInfoHandler())
	rg.GET("/:place_id/reviews", placeHandler.GetPlaceReviewsHandler())
	rg.POST("/review", placeHandler.WriteReviewHandler())
	rg.GET("/my-reviews", placeHandler.GetMyReviewsHandler())
	rg.DELETE("/review/:review_id", placeHandler.DeleteReviewHandler())
	rg.PATCH("review/:review_id", placeHandler.ModifyReviewHandler())
}

func registerPerformanceRoutes(rg *gin.RouterGroup, performanceHandler *handler.PerformanceHandler) {
	rg.GET("", performanceHandler.GetPerformanceShortList())  //목록조회
	rg.GET("/:id", performanceHandler.GetPerformanceDetail()) //공연 상세 조회

	rg.GET("/facility", performanceHandler.GetFacilityList())       //공연 시설 목록 조회
	rg.GET("/facility/:id", performanceHandler.GetFacilityDetail()) //공연 시설 상세 조회

	rg.GET("/top", performanceHandler.GetTopPerformances()) //topN 공연 조회
	rg.GET("/recommendation", performanceHandler.AiRecommendation())
	rg.GET("/near", performanceHandler.GetCommingPerformances())

	// rg.POST("/view", performanceHandler.IncrementPerformanceView())
}

func registerChatRoutes(rg *gin.RouterGroup, chatRoomHandler *handler.ChatRoomHandler) {
	rg.POST("/room", chatRoomHandler.CreateChatRoom())
	rg.POST("/join-room", chatRoomHandler.JoinChatRoom())
}
