package chat

import "github.com/SOMTHING-ITPL/ITPL-server/aws/dynamo"

func newChatRoomManager(db *dynamo.TableBasics) *RoomManager {
	return &RoomManager{
		rooms: make(map[uint]*Hub),
	}
}

func (m *RoomManager) GetOrCreate(roomID uint) *Hub {
	m.mu.RLock()
	hub, ok := m.rooms[roomID]
	m.mu.RUnlock()

	if ok {
		return hub
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.rooms[roomID]; ok {
		return hub
	}
	//double check

	hub = &Hub{
		roomID:     roomID,
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message, 256),
		closeCh:    make(chan struct{}),
		DynamoDB:   m.DynamoDB,
	}
	m.rooms[roomID] = hub
	go hub.Run()
	return hub

}

func (m *RoomManager) DeleteIfEmpty(roomID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.rooms[roomID]; ok {
		if len(hub.clients) == 0 {
			close(hub.closeCh) // chann 로 close
			delete(m.rooms, roomID)
		}
	}
}
