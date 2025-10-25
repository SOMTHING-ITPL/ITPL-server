package chat

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/google/uuid"
)

func BuildMessage(bucketBasics s3.BucketBasics, senderID, roomID uint, text string) Message {
	now := time.Now().UTC()
	sk := now.Format(time.RFC3339Nano) + "#" + uuid.NewString()
	return Message{
		MessageSK: sk,
		SenderID:  senderID,
		RoomID:    roomID,
		Timestamp: now,
		Text:      text,
	}
}
