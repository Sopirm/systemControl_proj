package models

type Role int

const (
	RoleObserver Role = iota
	RoleEngineer
	RoleManager
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}
