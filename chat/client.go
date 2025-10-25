package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client 로부터 모든 읽기 담당 단일 고루틴. 수신 메세지는 hub로 넘김.
func (c *Client) ReadMessages() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage() //MessageType이 TextMessage일 경우에만 사용함. Client 측에서 Pong 프레임을 보내야함.
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := BuildMessage(c.id, c.roomID, string(message))

		c.hub.broadcast <- msg
	}
}

func (c *Client) WriteMessages() {
	ticker := time.NewTicker(pingPeriod) //Client 측에서 pong 응답 줘야 함.
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok { //send channel 이 닫혔을 경우
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage) //websocket.TextMessage 는 json 형식으로 보내는 것.
			if err != nil {
				return
			}

			// Convert Message struct to JSON bytes
			messageBytes, err := json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(messageBytes)

			n := len(c.send) //send channel queue 길이
			for i := 0; i < n; i++ {
				w.Write(newline)
				queuedMessage := <-c.send
				queuedMessageBytes, err := json.Marshal(queuedMessage)
				if err != nil {
					continue
				}
				w.Write(queuedMessageBytes)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C: //ping 주기
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil { //ping 메세지 보내기
				return
			}
		}
	}
}
