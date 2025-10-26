package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

func NewChatRoomHandler(chatRoomRepo *chat.ChatRoomRepository, userRepo *user.Repository, bucketBasics *s3.BucketBasics, basics *dynamo.TableBasics, rm *chat.RoomManager) *ChatRoomHandler {
	return &ChatRoomHandler{
		chatRoomRepository: chatRoomRepo,
		userRepository:     userRepo,
		bucketBasics:       bucketBasics,
		tableBasics:        basics,
		chatRoomManager:    rm,
	}
}

// GET
// only title search
func (h *ChatRoomHandler) GetChatRoomsByTitle() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		rooms, err := h.chatRoomRepository.SearchChatRoomsByTitle(title)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		var response []ChatRoomInfoResponse
		for _, room := range rooms {
			roomInfo, err := ToChatRoomInfoResponse(h.bucketBasics.AwsConfig, h.bucketBasics.BucketName, room)
			if err != nil {
				log.Printf("Failed to get chat room info: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room info"})
				return
			}
			response = append(response, roomInfo)
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    response,
		})
	}
}

// GET
// search by coordinates, title and performance day
func (h *ChatRoomHandler) GetChatRoomsByCoordinate() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		ArrivalLongitude := c.Query("arrival_longitude")
		ArrivalLatitude := c.Query("arrival_latitude")

		arrivalLongitude, err := strconv.ParseFloat(ArrivalLongitude, 64)
		if err != nil {
			log.Printf("Failed to parse arrival longitude: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid arrival_longitude"})
			return
		}
		arrivalLatitude, err := strconv.ParseFloat(ArrivalLatitude, 64)
		if err != nil {
			log.Printf("Failed to parse arrival latitude: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid arrival_latitude"})
			return
		}

		DepartureLongitude := c.Query("departure_longitude")
		DepartureLatitude := c.Query("departure_latitude")

		departureLongitude, err := strconv.ParseFloat(DepartureLongitude, 64)
		if err != nil {
			log.Printf("Failed to parse departure longitude: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departure_longitude"})
			return
		}

		departureLatitude, err := strconv.ParseFloat(DepartureLatitude, 64)
		if err != nil {
			log.Printf("Failed to parse departure latitude: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid departure_latitude"})
			return
		}

		performanceDay := c.Query("performance_day")
		if ArrivalLongitude == "" || ArrivalLatitude == "" || DepartureLongitude == "" || DepartureLatitude == "" || performanceDay == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters"})
			return
		}
		perfDay, err := strconv.ParseInt(performanceDay, 10, 64)
		if err != nil {
			log.Printf("Failed to parse performance day: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid performance_day"})
			return
		}
		rooms, err := h.chatRoomRepository.GetChatRoomsByCoordinate(title, perfDay, departureLatitude, departureLongitude, arrivalLatitude, arrivalLongitude)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var response []ChatRoomInfoResponse
		for _, room := range rooms {
			roomInfo, err := ToChatRoomInfoResponse(h.bucketBasics.AwsConfig, h.bucketBasics.BucketName, room)
			if err != nil {
				log.Printf("Failed to get chat room info: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room info"})
				return
			}
			response = append(response, roomInfo)
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    response,
		})
	}
}

// GET
func (h *ChatRoomHandler) GetChatRoomMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}
		rid, err := strconv.ParseUint(roomID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
			return
		}
		chatRoom, err := h.chatRoomRepository.GetChatRoomById(uint(rid))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		members, err := h.chatRoomRepository.GetMembers(chatRoom)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		var response []ChatRoomMemberResponse
		for _, member := range members {
			memberInfo, err := ToChatRoomMemberInfoResponse(h.bucketBasics.AwsConfig, h.bucketBasics.BucketName, h.chatRoomRepository.DB, member.UserID)
			if err != nil {
				log.Printf("Failed to get chat room member info: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room member info"})
				return
			}
			response = append(response, memberInfo)
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: fmt.Sprintf("members of room %d", rid),
			Data:    response,
		})
	}
}

// GET
func (h *ChatRoomHandler) GetMyChatRooms() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		chatRooms, err := h.chatRoomRepository.GetMyChatRooms(userID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		var response []ChatRoomInfoResponse
		for _, room := range chatRooms {
			roomInfo, err := ToChatRoomInfoResponse(h.bucketBasics.AwsConfig, h.bucketBasics.BucketName, room)
			if err != nil {
				log.Printf("Failed to get chat room info: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room info"})
				return
			}
			response = append(response, roomInfo)
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "My chat rooms",
			Data:    response,
		})
	}
}

// POST
func (h *ChatRoomHandler) JoinChatRoom() gin.HandlerFunc {
	type request struct {
		RoomID uint `json:"room_id"`
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
		if err := h.chatRoomRepository.AddUserToChatRoom(h.userRepository, userID, req.RoomID); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// POST
func (h *ChatRoomHandler) CreateChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		params, err := chat.LoadParams(c) // chat/util.go
		if err != nil {
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

		info := chat.ChatRoomInfo{
			Title:              params.Title,
			ImgKey:             imageKey,
			PerformanceDay:     params.PerformanceDay,
			MaxMembers:         params.MaxMembers,
			DepartureLatitude:  params.DepartureLatitude,
			DepartureLongitude: params.DepartureLongitude,
			ArrivalLatitude:    params.ArrivalLatitude,
			ArrivalLongitude:   params.ArrivalLongitude,
			DepartureName:      params.DepartureName,
			ArrivalName:        params.ArrivalName,
		}

		err = h.chatRoomRepository.CreateChatRoom(h.userRepository, info, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// POST
func (h *ChatRoomHandler) LeaveChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}
		rid, err := strconv.ParseUint(roomID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		if err := h.chatRoomRepository.DeleteChatRoomMember(userID, uint(rid)); err != nil {
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
		mmap, err := h.tableBasics.GetItemsByPartitionKey(c, "room_id", &types.AttributeValueMemberN{Value: roomID})
		if err != nil {
			log.Printf("Failed to get items from DynamoDB : %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		messages, err := chat.MapToMessage(mmap)
		if err != nil {
			log.Printf("map to Message struct Convert Error %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		var response []ChatMessageResponse
		for _, message := range messages {
			res, err := ToChatMessageResponse(h.bucketBasics.AwsConfig, h.bucketBasics.BucketName, h.chatRoomRepository.DB, message)
			if err != nil {
				log.Printf("Failed to get chat room member info: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room member info"})
				return
			}
			response = append(response, res)
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    response,
		})
	}
}

func (h *ChatRoomHandler) ConnectToChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("room_id")
		log.Printf("[WS DEBUG] ConnectToChatRoom called with room_id: %s", roomID)

		if roomID == "" {
			log.Printf("[WS DEBUG] Missing room_id parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}
		rid, err := strconv.ParseUint(roomID, 10, 64)
		if err != nil {
			log.Printf("[WS DEBUG] Failed to parse room_id %s: %v", roomID, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room_id"})
			return
		}
		uid, ok := c.Get("userID")
		if !ok {
			log.Printf("[WS DEBUG] No userID found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID, _ := uid.(uint)
		log.Printf("[WS DEBUG] User %d attempting to connect to room %d", userID, rid)

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("[WS DEBUG] WebSocket upgrade failed for user %d, room %d: %v", userID, rid, err)
			return
		}
		log.Printf("[WS DEBUG] WebSocket connection established for user %d, room %d", userID, rid)

		chat.ServeWs(rid, userID, conn, h.chatRoomManager)

	}
}
