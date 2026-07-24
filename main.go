package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"context"
	"Http-Server/server"
	"Http-Server/database/bolt"
	"github.com/gorilla/mux"
)

func main() {
	//Http Variables
	address := ":9090"
	mux := mux.NewRouter()

	ctx := context.Background()
	blt := bolt.New(ctx)

	//Server
	srvr := server.New(blt)
	mux.HandleFunc("/", srvr.HandleIndex)
	//Using the .Methods() function to specify the allowed HTTP methods
	// Otherwise gorilla/mux will confuse /user/{name} with /user/create
	mux.HandleFunc("/user/{name}", srvr.HandleUser).Methods("GET", "PATCH","DELETE")
	mux.HandleFunc("/user/create", srvr.HandleCreateUser).Methods("POST", "PUT")
	
	
	//Http Server
	server := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//CLI Notification
	fmt.Println("Start Server ", address)
	log.Fatal(server.ListenAndServe())
}