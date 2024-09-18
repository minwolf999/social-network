package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	middleware "social-network/Middleware"
	routes "social-network/Routes"
	utils "social-network/Utils"
)

func init() {
	args := os.Args
	if len(args) != 2 {
		return
	}

	if strings.ToLower(args[1]) == "--loaddata" || strings.ToLower(args[1]) == "-l" {
		db, err := utils.OpenDb("sqlite3", "./Database/Database.sqlite")
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()
		
		if err = utils.LoadData(db); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	fmt.Println("\033[96mServer started at: http://localhost:8080\033[0m")

	// We create a log file and redirect the stdout to the new file
	logFile, _ := os.Create("./Log/" + time.Now().Format("2006-01-02__15-04-05") + ".log")
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("Server started at: http://localhost:8080")

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
		Handler:      handler,
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		// We start the listening of the port
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Signal to capture a clean shutdown (SIGINT/SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We close correctly the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s\n", err)
	}

	// We reset the stdout to is normal status
	fmt.Println("Server exiting")
	log.Println("Server exiting")

}
