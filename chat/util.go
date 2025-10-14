package chat

import (
	"encoding/json"
)

func MapToMessage(m []map[string]any) ([]Message, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var messages []Message
	err = json.Unmarshal(bytes, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ChatRoomRepository) GetNumOfMembers(room *ChatRoom) (int, error) {
	var count int64
	err := r.DB.Model(&ChatRoomMember{}).
		Where("chat_room_id = ?", room.ID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *ChatRoomRepository) IsFull(room *ChatRoom) (bool, error) {
	numOfMembers, err := r.GetNumOfMembers(room)
	if err != nil {
		return false, err
	}
	if room.MaxMembers == numOfMembers {
		return true, nil
	}
	return false, nil
}
