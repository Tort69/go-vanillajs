package data

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"Allusion/logger"
	"Allusion/models"
	"Allusion/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewAccountRepository(db *pgxpool.Pool, log *logger.Logger) (*AccountRepository, error) {
	return &AccountRepository{
		db:     db,
		logger: log,
	}, nil
}

func (r *AccountRepository) ResetPassword(email string, currentPassword string, newPassword string) (bool, error) {

	var user models.User
	query := `
		SELECT email, password_hashed
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.Email,
		&user.PasswordHashed,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}

	if err != nil {
		r.logger.Error("Failed to query user for reset PASSWORD", err)
		return false, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHashed), []byte(currentPassword))
	if err != nil {
		r.logger.Error("Password mismatch for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash new password", err)
		return false, err
	}

	query = `
	UPDATE users
	SET password_hashed = $1
	WHERE email = $2
	RETURNING id
`

	var userID int
	err = r.db.QueryRow(context.Background(),
		query,
		string(hashedPassword),
		email,
	).Scan(&userID)

	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}

	r.logger.Info("Succsec reset password user with email:" + email)
	return true, nil
}

func (r *AccountRepository) ResendVerifyEmail(email string) (bool, error) {

	var user models.User
	var exists bool
	err := r.db.QueryRow(context.Background(), `
		SELECT EXISTS(
		SELECT 1
			FROM users
			WHERE email = $1 AND is_verified = FALSE)
	`, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing user", err)
		return false, err
	}
	if !exists {
		r.logger.Error(" Unable to find the user or mail is already verified", ErrUserAlreadyExists)
		return false, ErrUserAlreadyExists
	}

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		r.logger.Error("Failed to create verify token", err)
	}

	user.VerifyToken = sql.NullString{String: token, Valid: token != ""}
	user.TokenExpiresAt = sql.NullTime{
		Time:  time.Now().Add(24 * time.Hour),
		Valid: true,
	}

	query := `
		UPDATE users
		SET verify_token = $1,
    token_expires_at = $2
		WHERE email = $3 AND is_verified = FALSE `

	commandTag, err := r.db.Exec(context.Background(), query, user.VerifyToken, user.TokenExpiresAt, email)
	if err != nil {
		r.logger.Error("Failed to save the verification token in the database", err)
		return false, ErrVerifyMail
	}
	if commandTag.RowsAffected() != 1 {
		r.logger.Error("Failed to save the verification token in the database", err)
		return false, ErrVerifyMail
	}
	utils.SendVerificationEmail(email, token)
	return true, nil
}

func (r *AccountRepository) Register(name, email, password string) (bool, error) {
	// Validate basic requirements
	var user models.User
	if name == "" || email == "" || password == "" {
		r.logger.Error("Registration validation failed: missing required fields", nil)
		return false, ErrRegistrationValidation
	}

	// Check if user already exists
	var exists bool
	err := r.db.QueryRow(context.Background(), `
		SELECT EXISTS(
		SELECT 1
			FROM users
			WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing user", err)
		return false, err
	}
	if exists {
		r.logger.Error("User already exists with email: "+email, ErrUserAlreadyExists)
		return false, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", err)
		return false, err
	}

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		r.logger.Error("Failed to create verify token", err)
	}

	user.IsVerified = false
	user.VerifyToken = sql.NullString{String: token, Valid: token != ""}
	user.TokenExpiresAt = sql.NullTime{
		Time:  time.Now().Add(24 * time.Hour),
		Valid: true,
	}

	// Insert new user
	query := `
		INSERT INTO users (name, email, password_hashed, time_created, is_verified, verify_token, token_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var userID int
	err = r.db.QueryRow(context.Background(),
		query,
		name,
		email,
		string(hashedPassword),
		time.Now(),
		user.IsVerified,
		user.VerifyToken,
		user.TokenExpiresAt,
	).Scan(&userID)
	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}
	utils.SendVerificationEmail(email, token)

	return true, nil
}

func (r *AccountRepository) Authenticate(email string, password string) (bool, error) {
	if email == "" || password == "" {
		r.logger.Error("Authentication validation failed: missing credentials", nil)
		return false, ErrAuthenticationValidation
	}

	// Fetch user by email
	var user models.User
	query := `
		SELECT id, name, email, password_hashed, is_verified, verify_token, token_expires_at
		FROM users
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHashed,
		&user.IsVerified,
		&user.VerifyToken,
		&user.TokenExpiresAt,
	)
	if user.IsVerified {
	}
	if err == sql.ErrNoRows {

		r.logger.Error("User not found for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}
	if err != nil {
		r.logger.Error("Failed to query user for authentication", err)
		return false, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHashed), []byte(password))
	if err != nil {
		r.logger.Error("Password mismatch for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}

	//verify email send
	if !user.IsVerified {
		token, err := utils.GenerateVerificationToken()
		if err != nil {
			r.logger.Error("Failed to create verify token", err)
		}

		user.IsVerified = false
		user.VerifyToken = sql.NullString{String: token, Valid: token != ""}
		user.TokenExpiresAt = sql.NullTime{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		}

		updateQuery := `
		UPDATE users
		SET last_login = $1 AND is_verified = FALSE AND verify_token= $1 AND token_expires_at= $2
		WHERE id = $3
		`
		commandTag, err := r.db.Exec(context.Background(), updateQuery, time.Now(), user.VerifyToken, user.TokenExpiresAt)
		if err != nil {
			r.logger.Error("Failed to save the verification token in the database", err)
			return false, ErrVerifyMail
		}
		if commandTag.RowsAffected() != 1 {
			r.logger.Error("Failed to save the verification token in the database", err)
			return false, ErrVerifyMail
		}
		utils.SendVerificationEmail(email, token)
		return false, ErrUserNotСonfirmedMail
	}

	// Update last login time
	updateQuery := `
		UPDATE users
		SET last_login = $1
		WHERE id = $2
	`

	commandTag, err := r.db.Exec(context.Background(), updateQuery, time.Now(), user.ID)
	if err != nil {
		r.logger.Error("Failed to update last login", err)
		// Don't fail authentication just because last login update failed
	}
	if commandTag.RowsAffected() != 1 {
		r.logger.Error("Failed to update last login", err)
		return false, ErrUserNotСonfirmedMail
	}

	return true, nil
}

func (r *AccountRepository) DeleteAccount(email string) (bool, error) {

	query := `
	DELETE FROM users WHERE email=$1
	RETURNING id
	`
	var user models.User

	err := r.db.QueryRow(context.Background(), query,
		email).Scan(
		&user.Email,
	)

	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}

	return true, nil

}

func (r *AccountRepository) GetAccountDetails(email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user by email", err)
		return models.User{}, err
	}

	// Fetch favorites
	favoritesQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year,
		       m.overview, m.score, m.popularity, m.language,
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'favorite'
		ORDER BY um.time_added
	`
	favoriteRows, err := r.db.Query(context.Background(), favoritesQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user favorites", err)
		return user, err
	}
	defer favoriteRows.Close()

	for favoriteRows.Next() {
		var m models.Movie
		if err := favoriteRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan favorite movie row", err)
			return user, err
		}
		user.Favorites = append(user.Favorites, m)
	}

	// Fetch watchlist
	watchlistQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year,
		       m.overview, m.score, m.popularity, m.language,
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'watchlist'
		ORDER BY um.time_added
	`
	watchlistRows, err := r.db.Query(context.Background(), watchlistQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user watchlist", err)
		return user, err
	}
	defer watchlistRows.Close()

	for watchlistRows.Next() {
		var m models.Movie
		if err := watchlistRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan watchlist movie row", err)
			return user, err
		}
		user.Watchlist = append(user.Watchlist, m)
	}

	return user, nil
}

func (r *AccountRepository) SaveCollection(user models.User, movieID int, collection string, score *int) (bool, error) {

	// Validate inputs
	if movieID <= 0 {
		r.logger.Error("SaveCollection failed: invalid movie ID", nil)
		return false, errors.New("invalid movie ID")
	}
	if collection != "favorite" && collection != "watchlist" {
		r.logger.Error("SaveCollection failed: invalid collection type", nil)
		return false, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	// Get user ID from email
	var userID int
	err := r.db.QueryRow(context.Background(), `
		SELECT id
		FROM users
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, err
	}

	// Insert the new relationship
	// if score != nil && collection == "favorite" {
	query := `
		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added, user_score)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, movie_id)
		DO UPDATE SET relation_type = EXCLUDED.relation_type, user_score = $5;
	`
	// } else {
	// 	query := `
	// 		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added, score)
	// 		VALUES ($1, $2, $3, $4)
	// 		ON CONFLICT (user_id, movie_id)
	// 		DO UPDATE SET relation_type = EXCLUDED.relation_type;
	// 	`

	// }
	_, err = r.db.Exec(context.Background(), query, userID, movieID, collection, time.Now(), score)
	if err != nil {
		r.logger.Error("Failed to save movie to "+collection, err)
		return false, err
	}

	r.logger.Info("Successfully added movie " + string(movieID) + " to " + collection + " for user")
	return true, nil
}

func (r *AccountRepository) DeleteCollection(user models.User, movieID int, collection string) (bool, error) {
	// Validate inputs
	if movieID <= 0 {
		r.logger.Error("DeleteCollection failed: invalid movie ID", nil)
		return false, errors.New("invalid movie ID")
	}
	if collection != "favorite" && collection != "watchlist" {
		r.logger.Error("DeleteCollection failed: invalid collection type", nil)
		return false, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	// Get user ID from email
	var userID int
	err := r.db.QueryRow(context.Background(), `
		SELECT id
		FROM users
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, err
	}

	// Delete the new relationship
	query := `
		DELETE FROM user_movies
		WHERE user_id = $1 AND movie_id = $2 AND relation_type = $3;
	`
	commandTag, err := r.db.Exec(context.Background(), query, userID, movieID, collection)
	if err != nil {
		r.logger.Error("Failed to delete movie to "+collection, err)
		return false, err
	}
	if commandTag.RowsAffected() != 1 {
		r.logger.Error("Failed to save the verification token in the database", err)
		return false, ErrVerifyMail
	}

	r.logger.Info("Successfully delete movie " + strconv.Itoa(movieID) + " to " + collection + " for user")
	return true, nil
}

func (r *AccountRepository) VerifyEmail(token string) (bool, string, error) {

	var user models.User
	query := `
        SELECT email
        FROM users
        WHERE
            verify_token = $1 AND
            is_verified = $2 AND
            token_expires_at > NOW()`

	err := r.db.QueryRow(context.Background(), query, token, false).Scan(&user.Email)
	r.logger.Info(user.Email)

	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, "", ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, "", err
	}

	updateQuery := `
	UPDATE users
		SET
    	verify_token = NULL,
    	token_expires_at = NULL,
    	is_verified = TRUE
		WHERE email = $1;`

	commandTag, err := r.db.Exec(context.Background(), updateQuery, user.Email)
	if err != nil {
		r.logger.Error("Failed to update verify email", err)
		return false, "", ErrVerifyMail
	}
	if commandTag.RowsAffected() != 1 {
		r.logger.Error("Failed to update verify email", err)
		return false, "", ErrVerifyMail
	}

	return true, user.Email, nil

}

var (
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrVerifyMail               = errors.New("email unconfirmed, confirmation sent to your e-mail again")
	ErrUserNotСonfirmedMail     = errors.New("user mail is not confirmed")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
)
