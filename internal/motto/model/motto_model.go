package model

type Motto struct {
	ID           int           `json:"id"`
	ImageURL     string        `json:"image_url" db:"image_url"`
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Name   string `json:"name" db:"name"`
	LangID int    `json:"language_id" db:"language_id"`
}
