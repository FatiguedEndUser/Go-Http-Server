package server_operations

import (
	"Http-Server/database"
	"Http-Server/pages"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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
		got := s.db.Get(s.ctx, user.Name)
		if got != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf("User already exists: %w", user.Name)))
			return
		}
		
		
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

	user := s.db.Get(s.ctx, name)
	if user == nil {
		writer.WriteHeader(http.StatusNotFound) //HTTP 404
		return
	}

	switch request.Method {
		case http.MethodGet:
			log.Printf("Get User: %s", name)
			msg, err := json.Marshal(user)
			if err != nil {
				log.Print("Could not marshal user: " + err.Error())
				writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
				return
			}

			writer.Header().Add("Content-Type", "application/json")
			writer.Write(msg)

		//Partial update
		case http.MethodPatch:
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
				writer.WriteHeader(http.StatusBadRequest) //HTTP 400
				return
			}

			//Validation
			if user.Name == "" {
				log.Print("Name is required")
				writer.WriteHeader(http.StatusBadRequest) //HTTP 400
				return
			}

			log.Printf("Update User: %s", name)
			s.db.Update(s.ctx, user)
			msg, err := json.Marshal(user)
			if err != nil {
				log.Print("Could not update user: " + err.Error())
				writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
				return
			}
			writer.Header().Add("Content-Type", "application/json")
			writer.Write(msg)
		
		case http.MethodDelete:
			log.Printf("Delete User: %s", name)
			err := s.db.Delete(s.ctx, name)
			if err != nil {
				log.Print("Could not delete user: " + err.Error())
				writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
				return
			}
			
		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed) //HTTP 405
	}
}
