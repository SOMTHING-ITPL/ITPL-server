package chat

import (
	"context"
	"log"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (c *ChatRoomMember) BroadcastMessage(room *ChatRoom, message Message, db *dynamodb.Client, tableBasics dynamo.TableBasics) {
	// 1. Save to DynamoDB
	go func() {
		av, err := attributevalue.MarshalMap(message)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}

		err = tableBasics.AddItemToDB(context.Background(), av)
		if err != nil {
			log.Println("Failed to save message to DynamoDB:", err)
		}
	}()

	// 2. WebSocket 브로드캐스트
	for i := range room.Members {
		member := room.Members[i]

		// Skip sender
		if member.UserID == c.UserID {
			continue
		}

		go func(m *ChatRoomMember) {
			m.Lock()
			defer m.Unlock()
			if m.Conn != nil {
				if err := m.Conn.WriteJSON(message); err != nil {
					log.Printf("Failed to send message to user %d: %v\n", m.UserID, err)
				}
			}
		}(member)
	}
}
