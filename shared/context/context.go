package context

import (
	"context"
	"errors"
	"net/http"
)

const (
	userIDKey   = "user_id"
	usernameKey = "username"
	roleKey     = "role"
	tokenKey    = "token"
)

// use failure.go
// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUsernameFromContext retrieves the username from the request context.
func GetUsernameFromContext(r *http.Request) (string, error) {
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}

func GetRoleFromContext(r *http.Request) (string, error) {
	role, ok := r.Context().Value(roleKey).(string)
	if !ok {
		return "", errors.New("role not found in context")
	}
	return role, nil
}

func GetTokenFromContext(r *http.Request) (string, error) {
	token, ok := r.Context().Value(tokenKey).(string)
	if !ok {
		return "", errors.New("token not found in context")
	}
	return token, nil
}

// WithUserID adds the user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithUsername adds the username to the context.
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}
