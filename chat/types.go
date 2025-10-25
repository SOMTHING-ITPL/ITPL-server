package chat

import (
	"sync"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// ChatRoom repository
type ChatRoomRepository struct {
	DB *gorm.DB
}

// WebSocket related structs
type WSRoom struct {
	sync.Mutex             // to avoid race condition
	RoomID     uint        // room ID
	Members    []*WSMember // userID -> ChatRoomMember
}

type WSMember struct {
	*websocket.Conn
	sync.Mutex

	UserID uint
}

type RoomManager struct {
	Rooms map[uint]*WSRoom // roomID -> WSRoom
	sync.Mutex
}

// for CreateChatRoom()
type ChatRoomInfo struct {
	Title          string  `json:"title"`
	ImgKey         *string `json:"img_key,omitempty"`
	PerformanceDay int64   `json:"performance_day"`
	MaxMembers     int     `json:"max_members"`
	DepartureCoord Region  `json:"departure_coord"`
	ArrivalCoord   Region  `json:"arrival_coord"`
	DepartureName  string  `json:"departure_name"`
	ArrivalName    string  `json:"arrival_name"`
}

// ChatRoom, ChatRoomMember : gorm mode
type Region struct {
	MapX float64 `json:"map_x"`
	MapY float64 `json:"map_y"`
}

type ChatRoom struct {
	gorm.Model
	Members        []*ChatRoomMember `json:"members" gorm:"foreignKey:ChatRoomID"`
	DepartureCoord Region            `json:"departure_coord"`
	ArrivalCoord   Region            `json:"arrival_coord"`

	DepartureName string `json:"departure_name"`
	ArrivalName   string `json:"arrival_name"`

	Title          string  `json:"title" gorm:"column:title"`
	ImageKey       *string `json:"image_key,omitempty" gorm:"column:image_key"`
	PerformanceDay int64   `json:"performance_day"`
	MaxMembers     int     `json:"max_members"`
}

type ChatRoomMember struct {
	User user.User `gorm:"foreignKey:UserID"`

	ChatRoomID uint      `json:"chat_room_id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"primaryKey"`
	IsAdmin    bool      `json:"is_admin" gorm:"column:is_admin;default:false"`
	JoinedAt   time.Time `json:"joined_at" gorm:"column:joined_at"`
}

type Message struct {
	MessageSK string    `json:"message_sk" dynamodbav:"message_sk"` // Sort Key (timestamp#uuid)
	SenderID  uint      `json:"sender_id" dynamodbav:"sender_id"`
	RoomID    uint      `json:"room_id" dynamodbav:"room_id"`     // Partition Key
	Timestamp time.Time `json:"timestamp" dynamodbav:"timestamp"` // stored as string RFC3339 fromat as default
	Text      string    `json:"text" dynamodbav:"text"`
}
