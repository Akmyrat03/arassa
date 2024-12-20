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

func (r *ImageRepository) Delete(id int) error {
	query := `DELETE FROM title WHERE id=$1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepository) GetImagePathsByTitleID(id int) ([]string, error) {
	var imagePaths []string
	query := `SELECT image_path FROM images WHERE title_id=$1`
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var imagePath string
		if err := rows.Scan(&imagePath); err != nil {
			return nil, err
		}
		imagePaths = append(imagePaths, imagePath)
	}
	return imagePaths, nil
}

func (r *ImageRepository) GetAllImages(langID int) ([]model.Image, error) {
	var images []model.Image
	query := `
		SELECT 
			t.id AS title_id,
			tt.lang_id,
			tt.title,
			i.image_path
		FROM 
			title AS t
		LEFT JOIN 
			title_translate AS tt ON t.id = tt.title_id
		LEFT JOIN 
			images AS i ON t.id = i.title_id
		WHERE tt.lang_id=$1	
		ORDER BY 
			t.id ASC, tt.lang_id ASC
	`
	err := r.DB.Select(&images, query, langID)
	if err != nil {
		return nil, err
	}

	return images, nil
}
