package users

import (
	"errors"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type User struct {
	ID        string     `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Password  string     `db:"password" json:"password"`
	Role      string     `db:"role" json:"role"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CreatedBy string     `db:"created_by" json:"created_by"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	UpdatedBy string     `db:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"` // * handle null
	DeletedBy *string    `db:"deleted_by" json:"deleted_by"`
}

type UserView struct {
	Username       string  `db:"username"`
	Name           *string `db:"name"`
	Role           string  `db:"role"`
	Gender         *string `db:"gender"`
	DoB            *string `db:"dob"`
	Education      *string `db:"education"`
	City           *string `db:"city"`
	Province       *string `db:"province"`
	Address        *string `db:"address"`
	PhoneNumber    *string `db:"phone_number"`
	JobRole        *string `db:"job_role"`
	Status         *string `db:"status"`
	PlacementCity  *string `db:"placement"`
	DepartmentName *string `db:"department_name"`
}

type UserList struct {
	Data         []UserView `json:"data"`
	TotalData    int        `json:"totalData"`
	TotalPages   int        `json:"totalPages"`
	CurrentPage  int        `json:"currentPage"`
	NextPage     *int       `json:"nextPage"`
	PreviousPage *int       `json:"previousPage"`
}

type UserFilter struct {
	Name     string `db:"name" json:"name"`
	City     string `db:"city" json:"city"`
	Province string `db:"province" json:"province"`
	JobRole  string `db:"job_role" json:"job_role"`
	Status   string `db:"status" json:"status"`
}

type ProfileView struct {
	Name           *string `db:"name"`
	Role           string  `db:"role"`
	Gender         *string `db:"gender"`
	DoB            *string `db:"dob"`
	Education      *string `db:"education"`
	City           *string `db:"city"`
	Province       *string `db:"province"`
	Address        *string `db:"address"`
	PhoneNumber    *string `db:"phone_number"`
	JobRole        *string `db:"job_role"`
	Status         *string `db:"status"`
	PlacementCity  *string `db:"placement_city"`
	DepartmentName *string `db:"department_name"`
}

type UpdateProfile struct {
	Name        *string   `db:"name" json:"name"`
	Gender      *string   `db:"gender" json:"gender"`
	DoB         *string   `db:"dob" json:"dob"`
	Education   *string   `db:"education" json:"education"`
	Address     *string   `db:"address" json:"address"`
	City        *string   `db:"city" json:"city"`
	Province    *string   `db:"province" json:"province"`
	PhoneNumber *string   `db:"phone_number" json:"phone_number"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
