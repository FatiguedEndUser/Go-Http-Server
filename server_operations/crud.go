package server

import (
	"Http-Server/pages"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"context"
	"Http-Server/database"
	"github.com/gorilla/mux"
)

// Default Constructor
func New(ctx context.Context, db database.Database) *Server {
	return &Server{
		ctx: ctx,
		db: db,
	}
}

// Server is an HTTP server.
type Server struct {
	ctx context.Context
	db database.Database
}


// User info is the information that is stored per user
type UserInfo struct {
	email string
	age   int
}

// HandleIndex handles the index path ("/") and serves a welcome message.
func (s *Server) HandleIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "text/html")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write([]byte(pages.Index))
}

// HandleCreateUser handles the create user path ("/create") and creates a new user.
// Create Post->Put
func (s *Server) HandleCreateUser(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodPost, http.MethodPut:

		//Check that the input type is json
		if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			log.Print("Could not read request body: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
			return
		}
		defer request.Body.Close()

		//Unmarshal the request body into a User struct
		var user database.User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Print("Could not unmarshal request body: " + err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//Validate
		// User not empty
		// User must not already exist

		if user.Name == "" {
			log.Print("Empty Username")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		

		
		//Write to the database
		err = s.db.Create(s.ctx, user)
		if err != nil {
			log.Print("Could not create user: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// // HandleUser handles the '/user' request for getting user information, and Updating user information.
func (s *Server) HandleUser(writer http.ResponseWriter, request *http.Request) {
// 	//Fetch the name from the query string
// 	//Common among all methods
	params := mux.Vars(request)
	name := params["name"]

	switch request.Method {
	case http.MethodGet:
		log.Printf("Get User: %s", name)
		user := s.db.Get(s.ctx, name)
		if user == nil {
			writer.WriteHeader(http.StatusNotFound) //HTTP 404
			return
		}

		msg, err := json.Marshal(user)
		if err != nil {
			log.Print("Could not marshal user: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Write(msg)

// 	//Partial update
// 	case http.MethodPatch:
// 		//Check that the input type is json
// 		if contentType := request.Header.Get("Content-Type"); contentType != "application/json" {
// 			writer.WriteHeader(http.StatusUnsupportedMediaType)
// 			return
// 		}

// 		body, err := io.ReadAll(request.Body)
// 		if err != nil {
// 			log.Print("Could not read request body: " + err.Error())
// 			writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
// 			return
// 		}
// 		defer request.Body.Close()

// 		//Unmarshal the request body into a User struct
// 		var user User
// 		err = json.Unmarshal(body, &user)
// 		if err != nil {
// 			log.Print("Could not unmarshal request body: " + err.Error())
// 			writer.WriteHeader(http.StatusBadRequest) //HTTP 400
// 			return
// 		}

// 		log.Printf("Update User: %s", name)
		
// 		//Get users
// 		userinfo := s.users[name]
// 		if user.Age != 0 {
// 			userinfo.age = user.Age
// 		}
// 		if user.Email != "" {
// 			userinfo.email = user.Email
// 		}
// 		s.users[name] = userinfo
// 		return

// 	case http.MethodDelete:
// 		log.Printf("Delete User: %s", name)
// 		delete(s.users, name)
// 		return
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed) //HTTP 405
	}
}
