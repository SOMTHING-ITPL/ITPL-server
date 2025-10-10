package chat

import (
	"context"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"
	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"gorm.io/gorm"
)

func CreateChatRoom(r *user.Repository, info ChatRoomInfo, myId uint) (*ChatRoom, error) {
	me, err := r.GetById(myId)
	if err != nil {
		return nil, err
	}
	creater := &ChatRoomMember{
		UserID:   me.ID,
		JoinedAt: time.Now(),
		IsAdmin:  true, // The creator is the admin
		User:     me,
	}
	return &ChatRoom{
		Title:          info.Title,
		ImageKey:       info.ImgKey,
		Members:        []*ChatRoomMember{creater},
		PerformanceDay: info.PerformanceDay,
		MaxMembers:     info.MaxMembers,
		Departure:      info.Departure,
		Arrival:        info.Arrival,
	}, nil
}

func GetChatRoomById(db *gorm.DB, roomId uint) (*ChatRoom, error) {
	var room ChatRoom
	if err := db.First(&room, roomId).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func JoinChatRoom(r *user.Repository, db *gorm.DB, userId uint, roomId uint) error {
	chatRoom, err := GetChatRoomById(db, roomId)
	if err != nil {
		return err
	}
	newUser, err := r.GetById(userId)
	if err != nil {
		return err
	}
	newMember := &ChatRoomMember{
		UserID:   newUser.ID,
		JoinedAt: time.Now(),
		IsAdmin:  false, // Regular member
		User:     newUser,
	}
	chatRoom.Members = append(chatRoom.Members, newMember)
	if err := db.Save(chatRoom).Error; err != nil {
		return err
	}
	return nil
}

func GetNumOfMembers(db *gorm.DB, room ChatRoom) (int, error) {
	var count int64
	err := db.Model(&ChatRoomMember{}).
		Where("chat_room_id = ?", room.ID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func IsFull(db *gorm.DB, room ChatRoom) (bool, error) {
	numOfMembers, err := GetNumOfMembers(db, room)
	if err != nil {
		return false, err
	}
	if room.MaxMembers == numOfMembers {
		return true, nil
	}
	return false, nil
}

func GetChatRoomsByCoordinate(db *gorm.DB, text string, performanceDay int64, departure Region, arrival Region) ([]ChatRoom, error) {
	var rooms []ChatRoom

	query := db.Model(&ChatRoom{}).
		Where("title LIKE ?", "%"+text+"%").
		Where("departure_map_x BETWEEN ? AND ?", departure.MapX-0.1, departure.MapX+0.1).
		Where("departure_map_y BETWEEN ? AND ?", departure.MapY-0.1, departure.MapY+0.1).
		Where("arrival_map_x BETWEEN ? AND ?", arrival.MapX-0.1, arrival.MapX+0.1).
		Where("arrival_map_y BETWEEN ? AND ?", arrival.MapY-0.1, arrival.MapY+0.1).
		Where("performance_day = ?", performanceDay)

	if err := query.Find(&rooms).Error; err != nil {
		return nil, err
	}

	return rooms, nil
}

func GetMembers(db *gorm.DB, room ChatRoom) ([]ChatRoomMember, error) {
	var members []ChatRoomMember
	err := db.Where("chat_room_id = ?", room.ID).Preload("User").Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func LeaveChatRoom(db *gorm.DB, userId uint, roomId uint) error {
	return db.Where("chat_room_id = ? AND user_id = ?", roomId, userId).Delete(&ChatRoomMember{}).Error
}

func DeleteChatRoom(ctx context.Context, gormDB *gorm.DB, bucketBasics *s3.BucketBasics, tableBasics *dynamo.TableBasics, roomID uint) error {
	sroomID := strconv.FormatUint(uint64(roomID), 10)
	// Delete Images from S3
	mmsg, err := tableBasics.GetItemsByPartitionKey(ctx, "room_id", &types.AttributeValueMemberN{Value: sroomID})
	if err != nil {
		return err
	}

	msg, err := MapToMessage(mmsg)
	if err != nil {
		return err
	}

	for _, m := range msg {
		if m.ContentType == "image" {
			if err := s3.DeleteImage(bucketBasics.S3Client, bucketBasics.BucketName, *m.ImageKey); err != nil {
				return err
			}
		}
	}
	/* Delete messages from DynamoDB */
	if err := tableBasics.DeleteItemsByPartitionKey(ctx, "room_id", &types.AttributeValueMemberN{Value: sroomID}); err != nil {
		return err
	}

	// Delete chatroom from MySQL
	return gormDB.Delete(&ChatRoom{}, roomID).Error
}
