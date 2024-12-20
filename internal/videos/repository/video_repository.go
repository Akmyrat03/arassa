package repository

import (
	"arassachylyk/internal/videos/model"

	"github.com/jmoiron/sqlx"
)

type VideoRepository struct {
	DB *sqlx.DB
}

func NewVideoRepository(DB *sqlx.DB) *VideoRepository {
	return &VideoRepository{DB: DB}
}

func (r *VideoRepository) Upload(title model.Title) (int, error) {
	tx, err := r.DB.Beginx()
	if err != nil {
		return 0, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Create video title
	var videoTitleID int
	err = tx.QueryRow("INSERT INTO video_titles DEFAULT VALUES RETURNING id").Scan(&videoTitleID)
	if err != nil {
		return 0, err
	}

	// Add translations
	for _, translation := range title.Translations {
		_, err := tx.Exec(
			"INSERT INTO video_title_translation (video_title_id, lang_id, video_title) VALUES ($1, $2, $3)",
			videoTitleID, translation.LangID, translation.Title,
		)
		if err != nil {
			return 0, err
		}
	}

	// Add video paths
	for _, video := range title.Videos {
		_, err := tx.Exec(
			"INSERT INTO videos (video_title_id, video_path) VALUES ($1, $2)",
			videoTitleID, video,
		)
		if err != nil {
			return 0, err
		}
	}

	return videoTitleID, nil
}

func (r *VideoRepository) Delete(id int) error {
	query := `DELETE FROM video_titles WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *VideoRepository) GetVideoPathsByID(id int) ([]string, error) {
	var paths []string

	query := "SELECT video_path FROM videos WHERE video_title_id = $1"
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var videoPath string
		if err := rows.Scan(&videoPath); err != nil {
			return nil, err
		}
		paths = append(paths, videoPath)
	}
	return paths, nil
}

func (r *VideoRepository) GetAllVideos(langID int) ([]model.Video, error) {
	var videos []model.Video
	query := `
		SELECT 
			vt.id AS video_title_id,
			vtt.lang_id,
			vtt.video_title,
			v.video_path
		FROM 
			video_titles AS vt
		LEFT JOIN 
			video_title_translation AS vtt ON vt.id = vtt.video_title_id
		LEFT JOIN 
			videos AS v ON vt.id = v.video_title_id
		WHERE vtt.lang_id=$1	
		ORDER BY 
			vt.id ASC, vtt.lang_id ASC
	`
	err := r.DB.Select(&videos, query, langID)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
