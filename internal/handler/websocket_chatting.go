package handler

import (
	_ "log"
	"net/http"

	_ "github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/gorilla/websocket"
)

// var allowedOrigins = []string{
// 	os.Getenv("FRONTEND_URL"),
// }

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,

// 	CheckOrigin: func(r *http.Request) bool {
// 		origin := r.Header.Get("Origin")
// 		for _, allowed := range allowedOrigins {
// 			if origin == allowed {
// 				return true // 일치하면 연결 허용
// 			}
// 		}
// 		return false
// 	},
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}
