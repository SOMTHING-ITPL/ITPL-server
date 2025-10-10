package chat

import (
	"context"
	"log"
	"strconv"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetChatHistory(ctx context.Context, bucketBasics *s3.BucketBasics, tableBasics *dynamo.TableBasics, roomID uint) ([]Message, error) {
	mmap, err := tableBasics.GetItemsByPartitionKey(ctx, "room_id", &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(roomID), 10)})
	if err != nil {
		log.Printf("room_id %d does not exist", roomID)
		return nil, err
	}
	messages, err := MapToMessage(mmap)
	if err != nil {
		log.Printf("failed to convert map to []Message. check chat/utils.go/MapToMessage()")
		return nil, err
	}
	for _, msg := range messages {
		if msg.ContentType == "image" {
			url, err := s3.GetPresignURL(bucketBasics.AwsConfig, bucketBasics.BucketName, *msg.ImageKey)
			if err != nil {
				log.Printf("failed to get presigned url(in Chat history provider)")
				return nil, err
			}
			msg.Content = &url
		}
	}
	return messages, nil
}
