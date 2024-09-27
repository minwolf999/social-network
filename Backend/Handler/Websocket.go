package handler

import (
	"database/sql"
	"net/http"
)

func Websocket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
