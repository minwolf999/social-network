package handler

import (
	"database/sql"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Websocket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		JWT := r.PathValue("JWT")
		userId, err := utils.DecryptJWT(JWT, db)
		if err != nil {
			log.Printf("[%s] [Websocket] There is an error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[%s] [Websocket] %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Unlock()
		model.ConnectedWebSocket.Conn[userId] = conn
		model.ConnectedWebSocket.Mu.Lock()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[%s] [Websocket] %v", r.RemoteAddr, err)
				break
			}

			conn.WriteJSON("Hello world!")
		}

		model.ConnectedWebSocket.Mu.Unlock()
		delete(model.ConnectedWebSocket.Conn, userId)
		model.ConnectedWebSocket.Mu.Lock()
	}
}
