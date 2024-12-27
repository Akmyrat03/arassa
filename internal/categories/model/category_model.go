package model

type CategoryDAO struct {
	ID           int           `json:"id" db:"id"`
	Translations []Translation `json:"translations" binding:"required"`
}

type Translation struct {
	Name   string `json:"name" binding:"required"`
	LangID int    `json:"lang_id" binding:"required"`
}

type CategoryReq struct {
	Translations []Translation `json:"translations" binding:"required"`
}

type CategoryRes struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}
