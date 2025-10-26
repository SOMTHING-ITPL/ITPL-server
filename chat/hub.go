package chat

import (
	"context"
	"log"
)

// TODO : send channel 이 가득 찰 경우에 => 블로킹 정책
func (h *Hub) Run() {
	log.Printf("[WS DEBUG] Hub.Run started")
	for {
		select {
		case client := <-h.register:
			log.Printf("[WS DEBUG] Client %d registered with hub", client.id)
			h.clients[client.id] = client
			log.Printf("[WS DEBUG] Total clients in hub: %d", len(h.clients))
		case client := <-h.unregister:
			log.Printf("[WS DEBUG] Client %d unregistering from hub", client.id)
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
				log.Printf("[WS DEBUG] Client %d removed from hub, total clients: %d", client.id, len(h.clients))
			}
		case msg := <-h.broadcast: //MessageType
			//save to database
			go func() {
				err := h.DynamoDB.AddItemToDB(context.Background(), msg)
				if err != nil {
					log.Printf("[WS DEBUG] Failed to save message to DB: %v", err)
				} else {
					log.Printf("[WS DEBUG] Message saved to DB successfully")
				}
			}()

			//broadcast to clients
			log.Printf("[WS DEBUG] Broadcasting message to %d clients", len(h.clients))
			for clientID, client := range h.clients {
				select {
				case client.send <- msg:
					log.Printf("[WS DEBUG] Message sent to client %d", clientID)
				default:
					log.Printf("[WS DEBUG] Send channel full for client %d, removing client", clientID)
					close(client.send) //client.send가 제대로 작동하지 않을경우.
					delete(h.clients, client.id)
				}
			}
		case <-h.closeCh:
			log.Printf("[WS DEBUG] Hub closing, disconnecting %d clients", len(h.clients))
			for _, client := range h.clients {
				close(client.send)
				client.conn.Close()
			}
			return
		}
	}
}
