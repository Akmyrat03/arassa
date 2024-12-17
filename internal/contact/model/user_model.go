package model

type ContactMessage struct {
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Message     string `json:"message" db:"message"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
}
