package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evermos/boilerplate-go/internal/domain/users"
	context_helpers "github.com/evermos/boilerplate-go/shared/context"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService    users.UserService
	Authentication *middleware.Authentication
}

func ProvideUserHandler(service users.UserService, auth *middleware.Authentication) UserHandler {
	return UserHandler{
		UserService:    service,
		Authentication: auth,
	}
}

func (h *UserHandler) Router(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Use(h.Authentication.VerifyJWT)
		r.Get("/profiles", h.GetProfile)
		r.Patch("/profiles", h.UpdateProfile)
		r.Group(func(r chi.Router) {
			r.Use(h.Authentication.IsAdmin)
			r.Get("/users", h.ReadUser)
			r.Delete("/users/{uuid}", h.DeleteUserByID)
		})
	})
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	name := q.Get("name")
	city := q.Get("city")
	province := q.Get("province")
	jobRole := q.Get("jobRole")
	status := q.Get("status")
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("size"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 5
	}

	response, err := h.UserService.ReadUser(users.UserFilter{
		Name:     name,
		City:     city,
		Province: province,
		JobRole:  jobRole,
		Status:   status,
	},
		page,
		size)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	err := h.UserService.DeleteUserByID(uuid)
	if err != nil {
		if strings.Contains(err.Error(), "admin role") {
			http.Error(w, "Cannot delete user with admin role or user not found", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "User deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	uuid, err := context_helpers.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from context", http.StatusInternalServerError)
		return
	}
	profile, err := h.UserService.GetProfile(uuid)
	if err != nil {
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var update users.UpdateProfile
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedBy, err := context_helpers.GetUsernameFromContext(r)
	if err != nil {
		http.Error(w, "Failed to get username from context", http.StatusInternalServerError)
		return
	}
	update.UpdatedBy = updatedBy

	uuid, err := context_helpers.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "Failed to get user ID from context", http.StatusInternalServerError)
		return
	}

	if update.DoB != nil && *update.DoB != "" {
		_, err := time.Parse("2006-01-02", *update.DoB)
		if err != nil {
			http.Error(w, "dob must be in YYYY-MM-DD format", http.StatusBadRequest)
			return
		}
	} else if update.DoB != nil && *update.DoB == "" {
		http.Error(w, "dob cannot be an empty string", http.StatusBadRequest)
		return
	}

	_, err = h.UserService.UpdateProfile(uuid, &update)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Profile updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
