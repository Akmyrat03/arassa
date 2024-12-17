package model

type Video struct {
	ID         int    `db:"id"`
	VideoPath  string `db:"video_path"`
	UploadedAt string `db:"uploaded_at"`
}
