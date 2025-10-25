package handler

import (
	_ "log"
	"net/http"
	"os"

	_ "github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/gorilla/websocket"
)

var allowedOrigins = []string{
	os.Getenv("FRONTEND_URL"),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true // 일치하면 연결 허용
			}
		}
		return false
	},
}

func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CloseWebSocket(conn *websocket.Conn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

/* 수정 필요
func (h *ChatHandler) ChatWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := UpgradeToWebSocket(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	var initMsg struct {
		UserID uint `json:"user_id"`
		RoomID uint `json:"room_id"`
	}

	err = conn.ReadJSON(&initMsg)
	if err != nil {
		log.Println("Failed to read init message:", err)
		return
	}

	user := &chat.ChatRoomMember{
		UserID: initMsg.UserID,
		Conn:   conn,
	}
	chatRoom, err := chat.GetChatRoomById(h.database, initMsg.RoomID)
	if err != nil {
		log.Println("Failed to get chat room:", err)
		return
	}

	go func() {
		for {
			var msg chat.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Failed to read message:", err)
				break
			}
			user.BroadcastMessage(chatRoom, msg, h.tableBasics.DynamoDbClient, *h.tableBasics)
		}
	}()

}
*/
