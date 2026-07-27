// Package server_operations provides HTTP handlers for managing user resources.
// It connects the HTTP layer to a database.Database implementation, translating
// requests into CRUD operations and responses.
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

// Server holds dependencies shared across all HTTP handlers.
type Server struct {
	ctx context.Context
	db  database.Database
}

// New creates a Server instance wired to the given database implementation.
func New(ctx context.Context, db database.Database) *Server {
	return &Server{
		ctx: ctx,
		db:  db,
	}
}

// HandleIndex serves the dashboard HTML page at the "/" route.
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(pages.Index))
}

// HandleCreateUser handles POST and PUT requests to "/user/create".
// Expects a JSON body matching the database.User struct.
// Validates that:
//   - Content-Type is application/json
//   - User name is non-empty
//   - User does not already exist in the database
func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:

		// Enforce JSON content type
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Printf("could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Parse JSON into User struct
		var user database.User
		if err := json.Unmarshal(body, &user); err != nil {
			log.Printf("could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate: name must be non-empty
		if user.Name == "" {
			log.Print("empty username in create request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate: user must not already exist
		if existing := s.db.Get(s.ctx, user.Name); existing != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("user already exists: %s", user.Name)))
			return
		}

		// Persist to database
		if err := s.db.Create(s.ctx, user); err != nil {
			log.Printf("could not create user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Return created user as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUser handles GET, PATCH, and DELETE requests to "/user/{name}".
// The {name} path variable identifies the target user.
//
// GET    - Returns the user as JSON (200), or 404 if not found.
// PATCH  - Updates the user's email/age with a JSON body (200), or 404/400 on failure.
// DELETE - Removes the user from the database (204), or 500 on failure.
func (s *Server) HandleUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Fetch user first — common to all methods
	user := s.db.Get(s.ctx, name)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:
		log.Printf("get user: %s", name)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("could not marshal user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodPatch:
		// Enforce JSON content type
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Printf("could not read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Parse JSON into User struct
		var updated database.User
		if err := json.Unmarshal(body, &updated); err != nil {
			log.Printf("could not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate: name must be non-empty
		if updated.Name == "" {
			log.Print("name is required for update")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("update user: %s", name)

		// Persist update — capture both return values
		result, err := s.db.Update(s.ctx, updated)
		if err != nil {
			log.Printf("could not update user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	case http.MethodDelete:
		log.Printf("delete user: %s", name)

		if err := s.db.Delete(s.ctx, name); err != nil {
			log.Printf("could not delete user: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}