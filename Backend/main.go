package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	routes "social-network/Routes"
	utils "social-network/Utils"
)

func init() {
	utils.CreateDb()
}

func main() {
	fmt.Println("\033[96mServer started at: http://localhost:8080\033[0m")

	mux := http.NewServeMux()

	routes.Routes(mux)

	srv := &http.Server{
		Handler: mux,
		Addr:    "localhost:8080",

		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
