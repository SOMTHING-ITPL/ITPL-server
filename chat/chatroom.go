package chat

import (
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"gorm.io/gorm"
)

func CreateChatRoom(r *user.Repository, title string, myId uint, performanceDay int64) (*ChatRoom, error) {
	me, err := r.GetById(myId)
	if err != nil {
		return nil, err
	}
	creater := ChatRoomMember{
		UserID:   me.ID,
		JoinedAt: time.Now(),
		IsAdmin:  true, // The creator is the admin
		User:     me,
	}
	return &ChatRoom{
		Title:          title,
		Members:        []ChatRoomMember{creater},
		PerformanceDay: performanceDay,
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
	newMember := ChatRoomMember{
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
