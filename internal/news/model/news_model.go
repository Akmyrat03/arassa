package model

import "time"

type News struct {
	ID           int           `db:"id"`
	CategoryID   int           `json:"category_id" binding:"required" db:"category_id"`
	ImageURL     string        `json:"image" binding:"required" db:"image"`
	CreatedAt    time.Time     `json:"created_at" binding:"required" db:"created_at"`
	Translations []Translation `json:"translations" binding:"required"`
}

type Translation struct {
	LangID      int    `json:"lang_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
