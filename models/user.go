package models

import "time"

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PasswordHashed string `json:"password"`
	IsVerified     bool `json:"is_verified"`
	VerifyToken    string `json:"verify_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	Favorites []Movie
	Watchlist []Movie
}
