package server

import (
	"Http-Server/pages"
	"log"
	"net/http"
	"io"
	"encoding/json"
)

//Default Constructor
func New() *Server {
	return &Server{
		//Fixes nil error when inputing json
		users: make(map[string]UserInfo),
	}
}

//Server is an HTTP server.
type Server struct {
	users map[string]UserInfo //Key -> username
}

//User information Types
	//User is the JSON value thats sent as a response to a request
type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Age int `json:"age"`
}
	//User info is the information that is stored per user
type UserInfo struct {
	email string
	age int
}

//HandleIndex handles the index path ("/") and serves a welcome message.
func (s *Server) HandleIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "text/html")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write([]byte(pages.Index))
}

//HandleUser handles the user path ("/user") and serves user information.
func (s *Server) HandleUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(pages.UserInfo))
}

//HandleCreateUser handles the create user path ("/create-user") and creates a new user.
//Create Post->Put
func (s *Server) HandleCreateUser(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost, http.MethodPut:
		//Check that the JSON input is correct
		contentType := request.Header.Get("Content-Type")
		if contentType != "application/json" {
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
		var user User
		err = json.Unmarshal(body, user)
		if err := json.Unmarshal(body, &user); err != nil {
			log.Print("Could not unmarshal request body: " + err.Error())
			writer.WriteHeader(http.StatusBadRequest) //HTTP 400
			return
		}

		log.Print("User created: " + user.Name)
		s.users[user.Name] = UserInfo{
			email: user.Email,
			age: user.Age,
		}
		
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
