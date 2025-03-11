package secure

import "github.com/gorilla/mux"

func RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/refresh", refreshTokenHandler).Methods("POST")

	// Protected routes
	api := router.PathPrefix("/users").Subrouter()
	api.Use(Protect)
	api.HandleFunc("", GetAllUsersHandler).Methods("GET")
	api.HandleFunc("/{id}", GetUserHandler).Methods("GET")
	api.HandleFunc("", CreateUserHandler).Methods("POST")
	api.HandleFunc("/{id}", UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/{id}", DeleteUserHandler).Methods("DELETE")
}
