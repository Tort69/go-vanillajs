package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"Allusion/data"
	"Allusion/logger"
	"Allusion/models"
	"Allusion/token"

	"github.com/redis/go-redis/v9" // indirect
)

// Define request structure
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Define request structure
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	JWT     string `json:"jwt"`
}

type AccountHandler struct {
	storage data.AccountStorage
	logger  *logger.Logger
	rdb     *redis.Client
}

// Utility functions
func (h *AccountHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Custom-Header", "value")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *AccountHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		switch err {
		case data.ErrAuthenticationValidation, data.ErrUserAlreadyExists, data.ErrRegistrationValidation:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return true
		case data.ErrVerifyMail:
			http.Error(w, "Error verify Email", http.StatusForbidden)
			return true

		case data.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
			return true
		case data.ErrUserNotСonfirmedMail:
			http.Error(w, "User mail is not confirmed", http.StatusForbidden)
			return true
		default:
			h.logger.Error(context, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return true
		}
	}
	return false
}

func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {

	// Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode registration request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the user
	success, err := h.storage.Register(req.Name, req.Email, req.Password)
	if !success {
		h.handleStorageError(w, err, "Failed to register user")
		return
	}

	// Return success response
	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully registered user with email: " + req.Email)
	}
}

func (h *AccountHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode authentication request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	success, err := h.storage.Authenticate(req.Email, req.Password)
	if !success {

		h.handleStorageError(w, err, "Failed to register user")
		return
	}
	// Return success response
	response := AuthResponse{
		Success: success,
		Message: "User log in successfully",
		JWT:     token.CreateJWT(models.User{Email: req.Email}, *h.logger),
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully authenticated user with email: " + req.Email)
	}
}

func (h *AccountHandler) HandlerResendVerifyEmail(w http.ResponseWriter, r *http.Request) {

	// Parse request body
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	// Register the user
	success, err := h.storage.ResendVerifyEmail(email)
	if !success {
		h.handleStorageError(w, err, "Failed to register user")
		return
	}

	// Return success response
	response := AuthResponse{
		Success: success,
		Message: "The email has been resent",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully send mail user with email: " + email)
	}
}

func (h *AccountHandler) HandlerResetPassword(w http.ResponseWriter, r *http.Request) {

	// Parse request body
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode registration request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.storage.ResetPassword(email, req.CurrentPassword, req.NewPassword)
	if !success {
		h.handleStorageError(w, err, "Failed to register user")
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
	}
	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully reset Password user with email: " + email)
	}

}

func (h *AccountHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := token.ExtractJWTSecret(r, *h.logger)
		if !ok {
			http.Error(w, "Email not found in token", http.StatusUnauthorized)
			return
		}

		// Inject email into the request context
		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const cooldown = 60 * time.Second

func (h *AccountHandler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.logger.Error("Failed to decode authentication request", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Проверка времени последнего запроса
		lastSent, err := h.rdb.Get(r.Context(), "resend:"+req.Email).Time()
		if err == nil && time.Since(lastSent) < cooldown {
			http.Error(w, "Email not found in token", http.StatusTooManyRequests)
			return
		}
		// Обновление времени
		err = nil
		if err := h.rdb.Set(r.Context(), "resend:"+req.Email, time.Now(), cooldown).Err(); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "email", req.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.storage.DeleteAccount(email)

	if h.handleStorageError(w, err, "Failed to save to collection") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "Successfully delete account",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully delete account")
	}
}

func (h *AccountHandler) SaveToCollection(w http.ResponseWriter, r *http.Request) {
	type CollectionRequest struct {
		MovieID    int    `json:"movie_id"`
		Collection string `json:"collection"`
		Score      *int   `json:"score"`
	}

	var req CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.storage.SaveCollection(models.User{Email: email},
		req.MovieID, req.Collection, req.Score)
	if h.handleStorageError(w, err, "Failed to save to collection") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "Movie added to " + req.Collection + " successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully saved movie to " + req.Collection)
	}
}

func (h *AccountHandler) DeleteToCollection(w http.ResponseWriter, r *http.Request) {

	type CollectionRequest struct {
		MovieID    int    `json:"movie_id"`
		Collection string `json:"collection"`
	}

	var req CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.storage.DeleteCollection(models.User{Email: email},
		req.MovieID, req.Collection)
	if h.handleStorageError(w, err, "Failed to save to collection") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "Movie delete to " + req.Collection + " successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully delete movie to " + req.Collection)
	}
}

func (h *AccountHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.storage.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := h.writeJSONResponse(w, details.Favorites); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}

func (h *AccountHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.storage.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := h.writeJSONResponse(w, details.Watchlist); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}

func (h *AccountHandler) VerifyByEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	h.logger.Info(token)

	IsVerified, email, err := h.storage.VerifyEmail(token)
	if h.handleStorageError(w, err, "Failed to verified") {
		return
	}
	if h.writeJSONResponse(w, IsVerified) == nil {
		h.logger.Info("Successfully verify by email: " + email)
	}
}

func NewAccountHandler(storage data.AccountStorage, log *logger.Logger, rdb *redis.Client) *AccountHandler {
	return &AccountHandler{
		storage: storage,
		logger:  log,
		rdb:     rdb,
	}
}
