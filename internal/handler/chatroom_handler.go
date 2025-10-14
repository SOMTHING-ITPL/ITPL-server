package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

func NewChatRoomHandler(chatRoomRepo *chat.ChatRoomRepository, userRepo *user.Repository, bucketBasics *s3.BucketBasics, basics *dynamo.TableBasics) *ChatRoomHandler {
	return &ChatRoomHandler{
		chatRoomRepository: chatRoomRepo,
		userRepository:     userRepo,
		bucketBasics:       bucketBasics,
		tableBasics:        basics,
	}
}

// POST
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

		var imageKey *string
		file, err := c.FormFile("image")
		if err == nil {
			// 업로드 처리
			key, err := s3.UploadToS3(
				h.bucketBasics.S3Client,
				h.bucketBasics.BucketName,
				fmt.Sprintf("chat_image/%d", userId),
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

		departure := chat.Region{
			MapX: request.DepartureLongitude,
			MapY: request.DepartureLatitude,
		}

		arrival := chat.Region{
			MapX: request.ArrivalLongitude,
			MapY: request.ArrivalLatitude,
		}

		info := chat.ChatRoomInfo{
			Title:          request.Title,
			ImgKey:         imageKey,
			PerformanceDay: request.PerformanceDay,
			MaxMembers:     request.MaxMembers,
			Departure:      departure,
			Arrival:        arrival,
		}

		err = h.chatRoomRepository.CreateChatRoom(h.userRepository, info, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Created chatroom successfully"})
	}
}

func (h *ChatRoomHandler) JoinChatRoom() gin.HandlerFunc {
	type request struct {
		UserID uint `json:"user_id"`
		RoomID uint `json:"room_id"`
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.chatRoomRepository.JoinChatRoom(h.userRepository, req.UserID, req.RoomID); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// GET
func (h *ChatRoomHandler) GetHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}
		m, err := h.tableBasics.GetItemsByPartitionKey(c, "room_id", &types.AttributeValueMemberN{Value: roomID})
		if err != nil {
			log.Printf("Failed to get items from DynamoDB : %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		messages, err := chat.MapToMessage(m)
		if err != nil {
			log.Printf("map to Message struct Convert Error %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		// 응답 구조 수정 필요
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}
