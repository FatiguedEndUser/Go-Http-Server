package server

import (
	"Http-Server/pages"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Default Constructor
func New() *Server {
	return &Server{
		//Fixes nil error when inputing json
		users: make(map[string]UserInfo),
	}
}

// Server is an HTTP server.
type Server struct {
	users map[string]UserInfo //Key -> username
}

// User information Types
// User is the JSON value thats sent as a response to a request
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
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
		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Print("Could not unmarshal request body: " + err.Error())
			writer.WriteHeader(http.StatusBadRequest) //HTTP 400
			return
		}

		log.Print("User created: " + user.Name)
		s.users[user.Name] = UserInfo{
			email: user.Email,
			age:   user.Age,
		}

	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUser handles the '/user' request for getting user information, and Updating user information.
func (s *Server) HandleUser(writer http.ResponseWriter, request *http.Request) {
	//Fetch the name from the query string
	//Common among all methods
	params := mux.Vars(request)
	name := params["name"]

	u, ok := s.users[name]
	if !ok {
		writer.WriteHeader(http.StatusNotFound) //HTTP 404
		return
	}

	switch request.Method {
	case http.MethodGet:
		returnedUser := User{
			Name:  name,
			Email: u.email,
			Age:   u.age,
		}

		msg, err := json.Marshal(returnedUser)
		if err != nil {
			log.Print("Could not marshal user: " + err.Error())
			writer.WriteHeader(http.StatusInternalServerError) //HTTP 500
			return
		}
		log.Printf("Get User: %s", name)

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
		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Print("Could not unmarshal request body: " + err.Error())
			writer.WriteHeader(http.StatusBadRequest) //HTTP 400
			return
		}

		log.Printf("Update User: %s", name)
		
		//Get users
		userinfo := s.users[name]
		if user.Age != 0 {
			userinfo.age = user.Age
		}
		if user.Email != "" {
			userinfo.email = user.Email
		}
		s.users[name] = userinfo
		return
		
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed) //HTTP 405
	}
}
