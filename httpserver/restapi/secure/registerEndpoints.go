package secure

import "github.com/gorilla/mux"

func RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	router.HandleFunc("/refresh", refreshTokenHandler).Methods("POST", "OPTIONS")

	// Protected routes
	api := router.PathPrefix("/users").Subrouter()
	api.Use(Protect)
	api.HandleFunc("", GetAllUsersHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("/{id}", GetUserHandler).Methods("GET", "OPTIONS")
	api.HandleFunc("", CreateUserHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/{id}", UpdateUserHandler).Methods("PUT", "OPTIONS")
	api.HandleFunc("/{id}", DeleteUserHandler).Methods("DELETE", "OPTIONS")
}
