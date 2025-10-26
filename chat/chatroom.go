package chat

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// // roomManager.rooms[roomID]에 user 추가
// func (rm *RoomManager) JoinUserToRoom(u *WSMember, roomID uint) {
// 	rm.Lock()
// 	defer rm.Unlock()
// 	room, exists := rm.Rooms[roomID]
// 	if !exists {
// 		rm.Rooms[roomID] = &WSRoom{
// 			Members: make([]*WSMember, 0),
// 		}
// 		room = rm.Rooms[roomID]
// 	}
// 	room.Lock()
// 	defer room.Unlock()
// 	room.Members = append(room.Members, u)
// }

// // roomManager.rooms[roomID]에서 user 제거
// func (rm *RoomManager) LeaveUserFromRoom(u *WSMember, roomId uint) error {
// 	rm.Lock()
// 	defer rm.Unlock()
// 	room, exists := rm.Rooms[roomId]
// 	if !exists {
// 		return fmt.Errorf("room %d does not exist", roomId)
// 	}
// 	room.Lock()
// 	defer room.Unlock()
// 	for i, member := range room.Members {
// 		if member == u {
// 			room.Members = append(room.Members[:i], room.Members[i+1:]...)
// 			break
// 		}
// 	}
// 	if len(room.Members) == 0 {
// 		delete(rm.Rooms, roomId)
// 	}
// 	return nil
// }

func ServeWs(w http.ResponseWriter, r *http.Request, rm *RoomManager) {
	roomIDStr := r.URL.Query().Get("room")
	userIDStr := r.URL.Query().Get("uid")
	if roomIDStr == "" || userIDStr == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid room ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	hub := rm.GetOrCreate(uint(roomID))
	client := &Client{
		id:     uint(userID),
		roomID: uint(roomID),
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
	}

	hub.register <- client

	go client.ReadMessages()
	go client.WriteMessages()

	return
}
