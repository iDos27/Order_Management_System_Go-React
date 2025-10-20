package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // Hasło nie jest zwracane w JSON
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // Minimum 6 znaków
	Role     string `json:"role" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // Minimum 6 znaków
	Role     string `json:"role" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Pod middleware/auth.go
type CurrentUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func IsAdmin(role string) bool {
	return role == RoleAdmin
}

func IsEmployee(role string) bool {
	return role == RoleEmployee
}

func IsCustomer(role string) bool {
	return role == RoleCustomer
}

func HasRole(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

const (
	RoleAdmin    = "admin"    // Admin
	RoleEmployee = "employee" // Pracownik
	RoleCustomer = "customer" // Klient

)
