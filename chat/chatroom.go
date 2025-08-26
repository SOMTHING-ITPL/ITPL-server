package chat

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func CreateChatRoom(
	r *user.Repository,
	title string,
	myId uint,
	performanceDay int64,
	maxMembers int,
	departure Region,
	arrival Region,
) (*ChatRoom, error) {
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
		Title:          title,
		Members:        []*ChatRoomMember{creater},
		PerformanceDay: performanceDay,
		MaxMembers:     maxMembers,
		Departure:      departure,
		Arrival:        arrival,
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

func GetCurrentMembers(db *gorm.DB, room ChatRoom) (int, error) {
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
	numOfMembers, err := GetCurrentMembers(db, room)
	if err != nil {
		return false, err
	}
	if room.MaxMembers == numOfMembers {
		return true, nil
	}
	return false, nil
}

/*
func GetChatRoomsByRegion(
	db *gorm.DB,
	text string,
	performanceDay int64,
	departure Region,
	arrival Region,
) ([]ChatRoom, error) {
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
*/
