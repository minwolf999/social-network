package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Websocket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[%s] [Websocket] %v", r.RemoteAddr, err)
			return
		}

		for {
			conn.WriteJSON("Hello world!")
		}
	}
}
