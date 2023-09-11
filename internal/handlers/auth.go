package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	AuthService    auth.AuthService
	Authentication *middleware.Authentication
}

func ProvideAuthHandler(service auth.AuthService, auth *middleware.Authentication) AuthHandler {
	return AuthHandler{
		AuthService:    service,
		Authentication: auth,
	}
}

// Router sets up the router for this domain.
func (h *AuthHandler) Router(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(h.Authentication.VerifyJWT)
			r.Use(h.Authentication.IsAdmin)
			r.Post("/register", h.Register)
		})
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the username and password from the request body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password fields are required", http.StatusBadRequest)
		return
	}

	// Authenticate user and generate JWT token
	token, err := h.AuthService.Login(req.Username, req.Password)
	if err != nil {
		if err == auth.ErrNotFound {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else if err == auth.ErrUnauthorized {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the JWT token
	w.Header().Set("Content-Type", "application/json")
	// use response.go
	response := map[string]string{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Define the required struct for the request body
	// Coba di define di model (auth model)
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "username, password, and role fields are required", http.StatusBadRequest)
		return
	}
	// di define, jangan hardcode
	adminUser := r.Context().Value("username").(string)
	user := &auth.User{
		Username:  req.Username,
		Password:  req.Password,
		Role:      req.Role,
		CreatedBy: adminUser,
		UpdatedBy: adminUser,
	}

	err := h.AuthService.Register(user)
	if err != nil {
		if err == auth.ErrUserExist {
			http.Error(w, "Username is already exist", http.StatusConflict)
		} else {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User successfully registered",
	}
	json.NewEncoder(w).Encode(response)
}
