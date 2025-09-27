package chat

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
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

func BuildImageMessage(cfg s3.BucketBasics, senderID, roomID uint, head *multipart.FileHeader) (ImageMessage, error) {
	key, err := s3.UploadToS3(cfg.S3Client, cfg.BucketName, "chat_images", head)
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

func BuildMessage(bucketBasics s3.BucketBasics, contentType string, message any) (Message, error) {
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
			imageURL, err := s3.GetPresignURL(bucketBasics.AwsConfig, bucketBasics.BucketName, msg.ImageKey)
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
