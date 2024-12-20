package model

type Title struct {
	ID           int           `db:"id"`
	Translations []Translation `json:"translations"`
	Videos       []string      `json:"videos"`
}

type Translation struct {
	LangID int    `db:"lang_id"`
	Title  string `db:"title"`
}

type Video struct {
	VideoTitleID int    `db:"video_title_id" json:"video_title_id"`
	LangID       int    `db:"lang_id" json:"lang_id"`
	VideoTitle   string `db:"video_title" json:"video_title"`
	VideoPath    string `db:"video_path" json:"video_path"`
}
