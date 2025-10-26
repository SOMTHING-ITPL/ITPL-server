package chat

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
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

func LoadParams(c *gin.Context) (ChatRoomInfo, error) {
	var info ChatRoomInfo
	title := c.PostForm("title")
	performanceDay, err := strconv.ParseInt(c.PostForm("performance_day"), 10, 64)
	if err != nil {
		return info, err
	}
	maxMembers, err := strconv.ParseInt(c.PostForm("max_members"), 10, 32)
	if err != nil {
		return info, err
	}
	maxMembersi := int(maxMembers)
	departureLongitude, err := strconv.ParseFloat(c.PostForm("departure_longitude"), 64)
	if err != nil {
		return info, err
	}
	departureLatitude, err := strconv.ParseFloat(c.PostForm("departure_latitude"), 64)
	if err != nil {
		return info, err
	}
	arrivalLongitude, err := strconv.ParseFloat(c.PostForm("arrival_longitude"), 64)
	if err != nil {
		return info, err
	}
	arrivalLatitude, err := strconv.ParseFloat(c.PostForm("arrival_latitude"), 64)
	if err != nil {
		return info, err
	}
	departureName := c.PostForm("departure_name")
	arrivalName := c.PostForm("arrival_name")

	info = ChatRoomInfo{
		Title:          title,
		PerformanceDay: performanceDay,
		MaxMembers:     maxMembersi,
		DepartureCoord: Region{
			MapX: departureLongitude,
			MapY: departureLatitude,
		},
		ArrivalCoord: Region{
			MapX: arrivalLongitude,
			MapY: arrivalLatitude,
		},
		DepartureName: departureName,
		ArrivalName:   arrivalName,
	}
	return info, nil
}
