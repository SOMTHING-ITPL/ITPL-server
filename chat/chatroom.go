package chat

import "fmt"

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[uint]*ChatRoom),
	}
}

func (rm *RoomManager) JoinUserToRoom(u *ChatRoomMember, roomId uint) error {
	rm.Lock()
	defer rm.Unlock()
	room, exists := rm.Rooms[roomId]
	if !exists {
		rm.Rooms[roomId] = &ChatRoom{
			Members: make([]*ChatRoomMember, 0),
		}
		room = rm.Rooms[roomId]
	}
	room.Lock()
	defer room.Unlock()
	room.Members = append(room.Members, u)
	return nil
}

func DeleteChatRoomMember(room *ChatRoom, userID uint) error {
	room.Lock()
	defer room.Unlock()
	for i, member := range room.Members {
		if member.UserID == userID {
			// Remove member from slice
			room.Members = append(room.Members[:i], room.Members[i+1:]...)
			break
		}
		if i == len(room.Members)-1 {
			return fmt.Errorf("Member with userID %d not found in room %d", userID, room.ID)
		}
	}
	return nil
}

func (rm *RoomManager) DeleteMember(member *ChatRoomMember, roomId uint) error {
	rm.Lock()
	defer rm.Unlock()
	room, exists := rm.Rooms[roomId]
	if !exists {
		return nil
	}
	room.Lock()
	defer room.Unlock()
	if err := DeleteChatRoomMember(room, member.UserID); err != nil {
		return err
	}
	return nil
}
