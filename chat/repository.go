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

func NewChatRoomRepository(db *gorm.DB) *ChatRoomRepository {
	return &ChatRoomRepository{DB: db}
}

func (r *ChatRoomRepository) CreateChatRoom(userRepo *user.Repository, info ChatRoomInfo, myID uint) error {
	me, err := userRepo.GetById(myID)
	if err != nil {
		return err
	}

	creater := &ChatRoomMember{
		UserID:   me.ID,
		JoinedAt: time.Now(),
		IsAdmin:  true, // The creator is the admin
		User:     me,
	}

	newChatRoom := &ChatRoom{
		Title:          info.Title,
		ImageKey:       info.ImgKey,
		Members:        []*ChatRoomMember{creater},
		PerformanceDay: info.PerformanceDay,
		MaxMembers:     info.MaxMembers,
		DepartureCoord: info.DepartureCoord,
		ArrivalCoord:   info.ArrivalCoord,
		DepartureName:  info.DepartureName,
		ArrivalName:    info.ArrivalName,
	}

	if err := r.DB.Create(newChatRoom).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChatRoomRepository) GetChatRoomById(roomId uint) (*ChatRoom, error) {
	var room ChatRoom
	if err := r.DB.First(&room, roomId).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *ChatRoomRepository) JoinChatRoom(userRepo *user.Repository, userId uint, roomId uint) error {
	chatRoom, err := r.GetChatRoomById(roomId)
	if err != nil {
		return err
	}
	newUser, err := userRepo.GetById(userId)
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
	if err := r.DB.Save(chatRoom).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChatRoomRepository) GetChatRoomsByCoordinate(text string, performanceDay int64, departure Region, arrival Region) ([]ChatRoom, error) {
	var rooms []ChatRoom

	query := r.DB.Model(&ChatRoom{}).
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

func (r *ChatRoomRepository) GetMembers(room *ChatRoom) ([]ChatRoomMember, error) {
	var members []ChatRoomMember
	err := r.DB.Where("chat_room_id = ?", room.ID).Preload("User").Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *ChatRoomRepository) DeleteChatRoomMember(userId uint, roomId uint) error {
	return r.DB.Where("chat_room_id = ? AND user_id = ?", roomId, userId).Delete(&ChatRoomMember{}).Error
}

func (r *ChatRoomRepository) DeleteChatRoom(ctx context.Context, bucketBasics *s3.BucketBasics, tableBasics *dynamo.TableBasics, roomID uint) error {
	sroomID := strconv.FormatUint(uint64(roomID), 10)
	/* Delete messages from DynamoDB */
	if err := tableBasics.DeleteItemsByPartitionKey(ctx, "room_id", &types.AttributeValueMemberN{Value: sroomID}); err != nil {
		return err
	}

	// Delete chatroom from MySQL
	return r.DB.Delete(&ChatRoom{}, roomID).Error
}
