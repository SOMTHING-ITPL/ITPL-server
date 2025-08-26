package chat

import (
	"sync"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Region struct {
	MapX float64 `json:"map_x"`
	MapY float64 `json:"map_y"`
}

type ChatRoom struct {
	gorm.Model
	Title          string            `json:"title" gorm:"column:title"`
	Members        []*ChatRoomMember `json:"members" gorm:"foreignKey:ChatRoomID"`
	PerformanceDay int64             `json:"performance_day"`
	MaxMembers     int               `json:"max_members"`
	Departure      Region            `json:"departure"`
	Arrival        Region            `json:"arrival"`
}

type ChatRoomMember struct {
	ChatRoomID uint `gorm:"primaryKey"`
	UserID     uint `gorm:"primaryKey"`
	IsAdmin    bool `gorm:"column:is_admin;default:false"`
	JoinedAt   time.Time
	User       user.User `gorm:"foreignKey:UserID"`
	Conn       *websocket.Conn
	Mu         sync.Mutex
}

type ChatMessage struct {
	Type      string    `json:"type"` // "text" or "image"
	SenderId  uint      `json:"sender"`
	Content   string    `json:"content"` // for text or image URL
	RoomId    uint      `json:"room_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
