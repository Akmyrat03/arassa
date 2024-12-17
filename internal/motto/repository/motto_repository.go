package repository

import (
	"arassachylyk/internal/motto/model"

	"github.com/jmoiron/sqlx"
)

const (
	Motto = "motto"
)

type MottoRepository struct {
	DB *sqlx.DB
}

func NewYearRepository(DB *sqlx.DB) *MottoRepository {
	return &MottoRepository{DB: DB}
}

func (r *MottoRepository) Create(motto model.Motto) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	var mottoId int
	query := `INSERT INTO motto (image_url) VALUES ($1) RETURNING id`
	err = tx.QueryRow(query, motto.ImageURL).Scan(&mottoId)
	if err != nil {
		return 0, err
	}

	for _, translation := range motto.Translations {
		query := `INSERT INTO motto_translate (motto_id, lang_id, name) VALUES ($1, $2, $3)`
		_, err := tx.Exec(query, mottoId, translation.LangID, translation.Name)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return mottoId, nil

}

func (r *MottoRepository) Delete(id int) error {
	query := `DELETE FROM motto WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MottoRepository) GetByID(id int) (model.Motto, error) {
	var motto model.Motto
	translations := []model.Translation{}
	query := `
		SELECT 
			m.id, m.image_url, mt.name, mt.lang_id 
		FROM 
			motto AS m 
		INNER JOIN 
			motto_translate AS mt ON m.id=mt.motto_id 
		WHERE m.id=$1`

	rows, err := r.DB.Query(query, id)
	if err != nil {
		return motto, err
	}
	defer rows.Close()

	for rows.Next() {
		var translation model.Translation
		err := rows.Scan(&motto.ID, &motto.ImageURL, &translation.Name, &translation.LangID)
		if err != nil {
			return motto, err
		}
		translations = append(translations, translation)
	}
	motto.Translations = translations
	return motto, nil
}
