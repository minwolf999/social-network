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

	middleware "social-network/Middleware"
	model "social-network/Model"
	routes "social-network/Routes"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	model.ConnectedWebSocket.Conn = make(map[string]*websocket.Conn)
	// tmp, err := bcrypt.GenerateFromPassword([]byte(model.SecretKey), 15)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// model.SecretKey = string(tmp)

	args := os.Args
	if len(args) != 2 {
		return
	}

	if strings.ToLower(args[1]) == "--loaddata" || strings.ToLower(args[1]) == "-l" {
		db, err := model.OpenDb("sqlite3", "./Database/Database.sqlite")
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()

		start := time.Now()
		if err = model.LoadData(db); err != nil {
			fmt.Println(err)
		}
		end := time.Now()
		fmt.Println(end.Sub(start))
	}
}

func main() {
	// We create a log file and redirect the stdout to the new file
	logFile, _ := os.Create("./Log/" + time.Now().Format("2006-01-02__15-04-05") + ".log")
	defer logFile.Close()

	log.SetOutput(logFile)

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
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Server listening on http://%s", srv.Addr)
		fmt.Printf("\033[96mServer started at: http://%s\033[0m\n", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Error starting TLS server: %v", err)
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
