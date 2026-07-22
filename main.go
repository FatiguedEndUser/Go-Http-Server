package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"Http-Server/server"
)

func main() {
	//Http Variables
	address := ":9090"
	mux := http.NewServeMux()

	//Server
	srvr := server.New()
	mux.HandleFunc("/", srvr.HandleIndex)
	mux.HandleFunc("/User", srvr.HandleUser)
	mux.HandleFunc("/User/Create-User", srvr.HandleCreateUser)
	
	
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