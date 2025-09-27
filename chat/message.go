package chat

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	aws_client "github.com/SOMTHING-ITPL/ITPL-server/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

func BuildTextMessage(senderID, roomID uint, text string) TextMessage {
	return TextMessage{
		SenderID:  senderID,
		Text:      text,
		RoomID:    roomID,
		Timestamp: time.Now(),
	}
}

func BuildImageMessage(cfg aws_client.BucketBasics, senderID, roomID uint, head *multipart.FileHeader) (ImageMessage, error) {
	key, err := aws_client.UploadToS3(cfg.S3Client, cfg.BucketName, "chat_images", head)
	if err != nil {
		return ImageMessage{}, err
	}
	return ImageMessage{
		SenderID:  senderID,
		RoomID:    roomID,
		Timestamp: time.Now(),
		ImageKey:  key,
	}, nil
}

func BuildMessage(bucketBasics aws_client.BucketBasics, contentType string, message any) (Message, error) {
	now := time.Now().UTC()
	sk := now.Format(time.RFC3339Nano) + "#" + uuid.NewString()

	if contentType == "text" {
		if msg, ok := message.(TextMessage); ok {
			content := msg.Text
			return Message{
				ContentType: "text",
				MessageSK:   sk,
				SenderID:    msg.SenderID,
				RoomID:      msg.RoomID,
				Timestamp:   msg.Timestamp,
				Content:     &content,
			}, nil
		}
	} else if contentType == "image" {
		if msg, ok := message.(ImageMessage); ok {
			imageURL, err := aws_client.GetPresignURL(bucketBasics.AwsConfig, bucketBasics.BucketName, msg.ImageKey)
			if err != nil {
				return Message{}, err
			}
			return Message{
				ContentType: "image",
				MessageSK:   sk,
				SenderID:    msg.SenderID,
				RoomID:      msg.RoomID,
				Timestamp:   msg.Timestamp,
				ImageURL:    &imageURL,
			}, nil
		}
	}

	err := fmt.Errorf("invalid content type. expected 'text' or 'image', got '%s'", contentType)
	return Message{}, err
}

func (c *ChatRoomMember) BroadcastMessage(room *ChatRoom, message Message, db *dynamodb.Client, tableName string) {
	go func() {
		av, err := attributevalue.MarshalMap(message)
		if err != nil {
			log.Println("Failed to marshal message:", err)
			return
		}

		_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      av,
		})
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
