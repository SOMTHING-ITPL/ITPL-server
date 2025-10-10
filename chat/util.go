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
