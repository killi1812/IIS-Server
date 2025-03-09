package secure

import "github.com/gorilla/mux"

func RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/refresh", refreshTokenHandler).Methods("POST")
}
