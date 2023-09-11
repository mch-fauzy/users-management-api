package auth

import (
	"errors"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrUserExist    = errors.New("username exist")
)

type Access struct {
	ID       string `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

type User struct {
	ID        string    `db:"id" json:"id"`
	ProfileID string    `db:"profile_id"`
	StatusID  string    `db:"status_id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}
