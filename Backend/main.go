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
	db, err := utils.OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	utils.CreateDb(db)
}

func main() {
	fmt.Println("\033[96mServer started at: http://localhost:8080\033[0m")

	// We launch the db
	mux := http.NewServeMux()

	// We set all the endpoints
	routes.Routes(mux)

	// We set the time out limit
	srv := &http.Server{
		Handler: mux,
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
