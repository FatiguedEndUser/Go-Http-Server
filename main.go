// Package main is the entry point for the HTTP user CRUD server.
// It initializes the BoltDB persistence layer, configures HTTP routing,
// and starts listening for requests on port 9090.
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"Http-Server/database/bolt"
	"Http-Server/server-operations"

	"github.com/gorilla/mux"
)

func main() {
	// Server configuration
	address := ":9090"

	// Initialize context for database startup
	ctx := context.Background()
	blt, err := bolt.New(ctx, "./data")
	if err != nil {
		log.Fatal("failed to start database: ", err)
	}
	defer func() {
		if err := blt.Close(ctx); err != nil {
			log.Print("warning: failed to close database: ", err)
		}
	}()

	// Timeout context for database operations
	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Set up router and server operations
	router := mux.NewRouter()
	srvr := server_operations.New(opCtx, blt)

	router.HandleFunc("/", srvr.HandleIndex)
	// Using .Methods() to specify allowed HTTP methods.
	// Without it, gorilla/mux would confuse /user/{name} with /user/create.
	router.HandleFunc("/user/{name}", srvr.HandleUser).Methods("GET", "PATCH", "DELETE")
	router.HandleFunc("/user/create", srvr.HandleCreateUser).Methods("POST", "PUT")

	// Configure HTTP server
	server := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("starting server on %s", address)
	log.Fatal(server.ListenAndServe())
}