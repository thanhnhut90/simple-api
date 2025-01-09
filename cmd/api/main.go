package main

import (
	"flag"
	"fmt"
	"github.com/thanhnhut90/simple-api/pkg/database"
	"github.com/thanhnhut90/simple-api/pkg/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func gracefulShutdown(server *http.Server, stopChan chan os.Signal) {
	<-stopChan
	fmt.Println("Shutting down gracefully...")
	if err := server.Shutdown(nil); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	fmt.Println("Server stopped.")
}

func main() {
	// Add a command-line flag to specify DB type
	dbType := flag.String("db", "postgres", "Database type (postgres or sqlite)")
	flag.Parse()

	// Initialize DB connection
	var err error
	db, err := database.InitDB(*dbType)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Conn.Close()

	// Create the strings table if it doesn't exist
	err = db.CreateTable("strings")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create the HTTPRouter object with DB
	httpRouter := &httprouter.HTTPRouter{
		DB: db,
	}

	// Setup routes
	httpRouter.SetupRoutes()

	// Graceful shutdown handling
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{Addr: ":8080"}
	go func() {
		fmt.Println("Server started on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	gracefulShutdown(server, stopChan)
}
