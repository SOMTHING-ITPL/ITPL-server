package chat

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

// TODO : send channel 이 가득 찰 경우에 => 블로킹 정책
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case msg := <-h.broadcast:
			//save to database
			go func() {
				av, err := attributevalue.MarshalMap(msg)
				if err != nil {
					log.Println("Failed to marshal message:", err)
					return
				}
				err = h.DynamoDB.AddItemToDB(context.Background(), av)
			}()

			for _, client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send) //client 강제 제거.
					delete(h.clients, client.id)
				}
			}
		case <-h.closeCh:
			for _, client := range h.clients {
				close(client.send)
				client.conn.Close()
			}
			return
		}
	}
}
