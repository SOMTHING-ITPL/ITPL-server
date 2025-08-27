package handler

import (
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewChatRoomHandler(db *gorm.DB, userRepo *user.Repository) *ChatRoomHandler {
	return &ChatRoomHandler{
		database:       db,
		userRepository: userRepo,
	}
}

func (h *ChatRoomHandler) CreateChatRoom() gin.HandlerFunc {
	type req struct {
		Title              string  `json:"title" binding:"required"`
		PerformanceDay     int64   `json:"performance_day" binding:"required"`
		MaxMembers         int     `json:"max_members" binding:"required"`
		DepartureLongitude float64 `json:"departure_longitude" binding:"required"`
		DepartureLatitude  float64 `json:"departure_latitude" binding:"required"`
		ArrivalLongitude   float64 `json:"arrival_longitude" binding:"required"`
		ArrivalLatitude    float64 `json:"arrival_latitude" binding:"required"`
	}
	return func(c *gin.Context) {
		var request req
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userId, _ := uid.(uint)

		departure := chat.Region{
			MapX: request.DepartureLongitude,
			MapY: request.DepartureLatitude,
		}

		arrival := chat.Region{
			MapX: request.ArrivalLongitude,
			MapY: request.ArrivalLatitude,
		}

		chatRoom, err := chat.CreateChatRoom(h.userRepository, request.Title, userId, request.PerformanceDay, request.MaxMembers, departure, arrival)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := h.database.Create(chatRoom).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Created chatroom successfully"})
	}
}

func (h *ChatRoomHandler) JoinChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
