package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

// 여기 단순히 roomid랑
func ServeWs(roomID uint64, userID uint, conn *websocket.Conn, rm *RoomManager) {
	log.Printf("[WS DEBUG] ServeWs called - roomID: %d, userID: %d", roomID, userID)

	hub := rm.GetOrCreate(uint(roomID)) //uint64 써야 할 것 같은데 급하니깐 keep dynamodb 랑 안 맞으면 큰일남.
	log.Printf("[WS DEBUG] Hub retrieved/created for room %d", roomID)

	client := &Client{
		id:     userID,
		roomID: uint(roomID),
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
	}
	log.Printf("[WS DEBUG] Client created for user %d in room %d", userID, roomID)

	hub.register <- client
	log.Printf("[WS DEBUG] Client %d registered with hub for room %d", userID, roomID)

	go client.ReadMessages()
	go client.WriteMessages()
	log.Printf("[WS DEBUG] Started read/write goroutines for user %d", userID)

	return
}
