package model

type Admin struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" binding:"required" db:"username"`
	Password string `json:"password" binding:"required" db:"password"`
}
