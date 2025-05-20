package models

import (
	"database/sql"
)

type User struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	PasswordHashed string         `json:"password"`
	IsVerified     bool           `json:"is_verified"`
	VerifyToken    sql.NullString `json:"verify_token"`
	TokenExpiresAt sql.NullTime      `json:"token_expires_at"`
	Favorites      []Movie
	Watchlist      []Movie
}
