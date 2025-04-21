package models

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PasswordHashed string `json:"password"`
	Favorites      []Movie
	Watchlist      []Movie
}
