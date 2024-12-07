package repository

import (
	"arassachylyk/internal/motto/model"
	"fmt"

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
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, language_id, image_url) VALUES ($1, $2, $3) RETURNING id", Motto)
	row := r.DB.QueryRow(query, motto.Name, motto.LanguageID, motto.ImageURL)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MottoRepository) GetByID(id int) (model.Motto, error) {
	var year model.Motto
	query := fmt.Sprintf("SELECT name, image_url FROM %v WHERE id=$1", Motto)
	err := r.DB.Get(&year, query, id)
	if err != nil {
		return model.Motto{}, err
	}
	return year, nil
}

func (r *MottoRepository) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id=$1", Motto)
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
