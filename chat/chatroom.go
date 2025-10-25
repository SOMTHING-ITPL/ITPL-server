package chat

import "fmt"

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[uint]*WSRoom),
	}
}

// roomManager.rooms[roomID]에 user 추가
func (rm *RoomManager) JoinUserToRoom(u *WSMember, roomID uint) {
	rm.Lock()
	defer rm.Unlock()
	room, exists := rm.Rooms[roomID]
	if !exists {
		rm.Rooms[roomID] = &WSRoom{
			Members: make([]*WSMember, 0),
		}
		room = rm.Rooms[roomID]
	}
	room.Lock()
	defer room.Unlock()
	room.Members = append(room.Members, u)
}

// roomManager.rooms[roomID]에서 user 제거
func (rm *RoomManager) LeaveUserFromRoom(u *WSMember, roomId uint) error {
	rm.Lock()
	defer rm.Unlock()
	room, exists := rm.Rooms[roomId]
	if !exists {
		return fmt.Errorf("room %d does not exist", roomId)
	}
	room.Lock()
	defer room.Unlock()
	for i, member := range room.Members {
		if member == u {
			room.Members = append(room.Members[:i], room.Members[i+1:]...)
			break
		}
	}
	if len(room.Members) == 0 {
		delete(rm.Rooms, roomId)
	}
	return nil
}
