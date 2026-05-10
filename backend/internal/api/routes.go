package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"car-lock-system/backend/internal/auth"
	"car-lock-system/backend/internal/db"
	"car-lock-system/backend/pkg/models"

	"github.com/gorilla/mux"
)

// MongoDB repositories
var (
	userRepo    *db.UserRepository
	bookingRepo *db.BookingRepository
)

// In-memory lock status (can be moved to MongoDB later)
var (
	mu       sync.RWMutex
	isLocked = false
)

// InitRepositories initializes MongoDB repositories
func InitRepositories() {
	userRepo = db.NewUserRepository()
	bookingRepo = db.NewBookingRepository()
}

// RegisterRoutes registers all API routes used by the backend. main.go
// expects this signature: RegisterRoutes(router *mux.Router)
func RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	// Auth
	api.HandleFunc("/register", RegisterHandler).Methods("POST")
	api.HandleFunc("/login", LoginHandler).Methods("POST")

	// User approval flow (admin approves/declines registrations)
	api.HandleFunc("/users/{username}/approve", ApproveUserHandler).Methods("POST")
	api.HandleFunc("/users/{username}/reject", RejectUserHandler).Methods("POST")
	api.HandleFunc("/users/pending", GetPendingUsersHandler).Methods("GET")

	// Lock control
	api.HandleFunc("/lock", LockHandler).Methods("POST")
	api.HandleFunc("/unlock", UnlockHandler).Methods("POST")
	api.HandleFunc("/status", StatusHandler).Methods("GET")

	// Booking flow (request -> admin approves/declines)
	api.HandleFunc("/bookings", CreateBookingHandler).Methods("POST")
	api.HandleFunc("/bookings/{id}", GetBookingHandler).Methods("GET")
	api.HandleFunc("/bookings/{id}/approve", ApproveBookingHandler).Methods("POST")
	api.HandleFunc("/bookings/{id}/decline", DeclineBookingHandler).Methods("POST")
	api.HandleFunc("/bookings/{id}/cancel", CancelBookingHandler).Methods("POST")
	api.HandleFunc("/bookings/{id}/reschedule", RescheduleBookingHandler).Methods("PATCH")
}

// helper: write JSON
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// RegisterHandler registers a new user with MongoDB. Passwords are hashed
// using bcrypt. New registrations start with pending_approval status
// and must be approved by admin before use.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	exists, err := userRepo.UserExists(ctx, req.Username)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
		return
	}
	if exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "username exists"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}

	// Create user document
	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Status:   "pending_approval",
	}

	err = userRepo.CreateUser(ctx, user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "pending_admin_approval", "username": req.Username})
}

// LoginHandler authenticates a user using MongoDB and bcrypt password comparison.
// Returns a placeholder JWT token (should be replaced with real JWT generation).
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get user from MongoDB
	user, err := userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Check user approval status
	if user.Status != "approved" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "user not approved by admin"})
		return
	}

	// Compare hashed password
	if !auth.ComparePassword(user.Password, req.Password) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	// Placeholder token - replace with real JWT in production
	writeJSON(w, http.StatusOK, map[string]string{"token": "placeholder-token", "username": req.Username})
}

// LockHandler locks the car. In the final system this should verify
// the caller is authorized and call into a service that updates MongoDB
// and notifies the car hardware if necessary.
func LockHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	isLocked = true
	mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"status": "locked"})
}

// UnlockHandler unlocks the car.
func UnlockHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	isLocked = false
	mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]string{"status": "unlocked"})
}

// StatusHandler returns current lock status
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	locked := isLocked
	mu.RUnlock()
	writeJSON(w, http.StatusOK, map[string]bool{"locked": locked})
}

// CreateBookingHandler allows a user to request a booking for the car.
// Bookings are stored in MongoDB with pending status and must be reviewed by admin.
func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User string `json:"user"`
		From string `json:"from"`
		To   string `json:"to"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking := &models.Booking{
		Username: req.User,
		From:     req.From,
		To:       req.To,
		Note:     req.Note,
		Status:   "pending",
	}

	err := bookingRepo.CreateBooking(ctx, booking)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create booking"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"booking_id": booking.ID.Hex(), "status": "pending"})
}

// GetBookingHandler returns the details of a single booking.
func GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
		return
	}

	writeJSON(w, http.StatusOK, booking)
}

// ApproveBookingHandler approves a pending booking (admin action).
// Updates booking status in MongoDB to "approved".
func ApproveBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify booking exists
	booking, err := bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
		return
	}

	// Update status
	err = bookingRepo.UpdateBookingStatus(ctx, id, "approved")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to approve booking"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"booking_id": id, "status": "approved", "username": booking.Username})
}

// DeclineBookingHandler declines a pending booking (admin action).
// Updates booking status in MongoDB to "declined".
func DeclineBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify booking exists
	booking, err := bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
		return
	}

	// Update status
	err = bookingRepo.UpdateBookingStatus(ctx, id, "declined")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to decline booking"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"booking_id": id, "status": "declined", "username": booking.Username})
}

// CancelBookingHandler allows a user to cancel a booking.
// Updates booking status in MongoDB to "cancelled".
func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
		return
	}

	err = bookingRepo.UpdateBookingStatus(ctx, id, "cancelled")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to cancel booking"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"booking_id": id, "status": "cancelled", "username": booking.Username})
}

// RescheduleBookingHandler allows a user to change the booking time or note.
func RescheduleBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
		Note string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	booking, err := bookingRepo.GetBookingByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
		return
	}

	if req.From == "" && req.To == "" && req.Note == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "nothing to update"})
		return
	}

	err = bookingRepo.UpdateBookingDetails(ctx, id, req.From, req.To, req.Note)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to reschedule booking"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"booking_id": id, "status": booking.Status, "username": booking.Username})
}

// ApproveUserHandler approves a pending user registration (admin action).
// Updates user status in MongoDB to "approved".
func ApproveUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update user status to approved
	err := userRepo.UpdateUserStatus(ctx, username, "approved")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to approve user"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"username": username, "status": "approved"})
}

// RejectUserHandler rejects a pending user registration (admin action).
// Updates user status in MongoDB to "rejected".
func RejectUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update user status to rejected
	err := userRepo.UpdateUserStatus(ctx, username, "rejected")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to reject user"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"username": username, "status": "rejected"})
}

// GetPendingUsersHandler returns a list of users pending admin approval.
func GetPendingUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := userRepo.GetUsersByStatus(ctx, "pending_approval")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get pending users"})
		return
	}

	// Convert to response format (exclude passwords)
	type UserResponse struct {
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pending_users": response,
		"count":         len(response),
	})
}
