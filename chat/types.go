package chat

import (
	"sync"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// for CreateChatRoom()
type ChatRoomInfo struct {
	Title          string  `json:"title"`
	ImgKey         *string `json:"img_key,omitempty"`
	PerformanceDay int64   `json:"performance_day"`
	MaxMembers     int     `json:"max_members"`
	Departure      Region  `json:"departure"`
	Arrival        Region  `json:"arrival"`
}

// ChatRoom, ChatRoomMember : gorm mode
type Region struct {
	MapX float64 `json:"map_x"`
	MapY float64 `json:"map_y"`
}

type ChatRoom struct {
	gorm.Model
	Members   []*ChatRoomMember `json:"members" gorm:"foreignKey:ChatRoomID"`
	Departure Region            `json:"departure"`
	Arrival   Region            `json:"arrival"`

	Title          string  `json:"title" gorm:"column:title"`
	ImageKey       *string `json:"image_key,omitempty" gorm:"column:image_key"`
	PerformanceDay int64   `json:"performance_day"`
	MaxMembers     int     `json:"max_members"`
}

type ChatRoomMember struct {
	User user.User `gorm:"foreignKey:UserID"`
	sync.Mutex
	*websocket.Conn

	ChatRoomID uint      `json:"chat_room_id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"primaryKey"`
	IsAdmin    bool      `json:"is_admin" gorm:"column:is_admin;default:false"`
	JoinedAt   time.Time `json:"joined_at" gorm:"column:joined_at"`
}

// Message : dynamodb model
type TextMessage struct {
	SenderID  uint      `json:"sender"`
	Text      string    `json:"text"`
	RoomID    uint      `json:"room_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type ImageMessage struct {
	SenderID  uint      `json:"sender_id"`
	RoomID    uint      `json:"room_id"`
	Timestamp time.Time `json:"timestamp"`
	ImageKey  string    `json:"image_key"`
}

type Message struct {
	MessageSK   string    `json:"message_sk" dynamodbav:"message_sk"`     // Sort Key (timestamp#uuid)
	ContentType string    `json:"content_type" dynamodbav:"content_type"` // "text" or "image"
	SenderID    uint      `json:"sender_id" dynamodbav:"sender_id"`
	RoomID      uint      `json:"room_id" dynamodbav:"room_id"`                         // Partition Key
	Timestamp   time.Time `json:"timestamp" dynamodbav:"timestamp"`                     // stored as string RFC3339 fromat as default
	Content     *string   `json:"content,omitempty" dynamodbav:"content,omitempty"`     // for text messages
	ImageURL    *string   `json:"image_url,omitempty" dynamodbav:"image_url,omitempty"` // for image messages
}
