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

type MottoResponse struct {
	MottoID  int    `db:"id" json:"motto_id"`
	LangID   int    `db:"lang_id" json:"lang_id"`
	Name     string `db:"name" json:"name"`
	ImageURL string `db:"image_url" json:"image_url"`
}
