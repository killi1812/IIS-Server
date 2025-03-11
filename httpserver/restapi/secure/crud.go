package secure

import (
	"encoding/json"
	"iis_server/apiq"
	"net/http"

	"github.com/gorilla/mux"
)

// GET /users - Get all users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetAllUsers())
}

// GET /users/{id} - Get user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, exists := GetUserByID(params["id"])
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// POST /users - Create user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user apiq.InstagramUsername
	json.NewDecoder(r.Body).Decode(&user)
	CreateUser(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// PUT /users/{id} - Update user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedUser apiq.InstagramUsername
	json.NewDecoder(r.Body).Decode(&updatedUser)
	if !UpdateUser(params["id"], updatedUser) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(updatedUser)
}

// DELETE /users/{id} - Delete user
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if !DeleteUser(params["id"]) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
