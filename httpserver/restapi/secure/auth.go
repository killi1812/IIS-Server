package secure

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TODO: see if better to use form or body
type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handler (generates tokens)
func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Dummy authentication (Replace with real authentication)
	if username != "admin" || password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := generateTokens(username)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

// Refresh token handler
func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse and verify refresh token
	token, err := jwt.ParseWithClaims(req.RefreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jefreshSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || claims.ExpiresAt.Unix() < time.Now().Unix() {
		http.Error(w, "Expired refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := generateTokens(claims.Username)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
}
