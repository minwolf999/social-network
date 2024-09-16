package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	middleware "social-network/Middleware"
	routes "social-network/Routes"
)

func main() {
	fmt.Println("\033[96mServer started at: http://localhost:8080\033[0m")

	// We launch the server
	mux := http.NewServeMux()

	// Enchaîner les middlewares
	handler := middleware.SetHeaderAccessControll(
		middleware.LookMethod(mux),
	)

	// We set all the endpoints
	routes.Routes(mux)

	// We set the time out limit
	srv := &http.Server{
		Handler: handler,
		Addr:    "localhost:8080",

		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	// We start the listening of the port
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
