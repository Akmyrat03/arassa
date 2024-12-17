package model

type Title struct {
	ID           int           `db:"id"`
	Translations []Translation `json:"translations"`
	Images       []string      `json:"images"`
}

type Translation struct {
	LangID int    `db:"lang_id"`
	Title  string `db:"title"`
}
