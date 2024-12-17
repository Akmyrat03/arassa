package repository

import (
	"arassachylyk/internal/images/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type ImageRepository struct {
	DB *sqlx.DB
}

func NewImageRepository(DB *sqlx.DB) *ImageRepository {
	return &ImageRepository{DB: DB}
}

func (r *ImageRepository) Create(title model.Title) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	var titleID int
	query := `INSERT INTO title DEFAULT VALUES RETURNING id`
	err = tx.QueryRow(query).Scan(&titleID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, translation := range title.Translations {
		query := `INSERT INTO title_translate (title_id, lang_id, title) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, titleID, translation.LangID, translation.Title)
		if err != nil {
			return 0, err
		}
	}

	imageQuery := `INSERT INTO images (title_id, image_path) VALUES ($1, $2)`
	for _, image := range title.Images {
		_, err := tx.Exec(imageQuery, titleID, image)
		if err != nil {
			log.Println("Error inserting images: ", err)
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction: ", err)
		return 0, err
	}

	return titleID, nil

}
