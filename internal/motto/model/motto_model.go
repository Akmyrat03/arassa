package model

type Motto struct {
	Name       string `json:"name" db:"name"`
	LanguageID int    `json:"language_id" db:"language_id"`
	ImageURL   string `json:"image_url" db:"image_url"`
}
