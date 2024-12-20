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

type Image struct {
	TitleID   int    `db:"title_id" json:"title_id"`
	LangID    int    `db:"lang_id" json:"lang_id"`
	Title     string `db:"title" json:"title"`
	ImagePath string `db:"image_path" json:"image_path"`
}
