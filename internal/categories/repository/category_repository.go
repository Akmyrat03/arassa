package repository

import (
	"arassachylyk/internal/categories/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	CATEGORIES = "categories"
)

type CategoryRepository struct {
	DB *sqlx.DB
}

func NewCategoryRepository(DB *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{DB: DB}
}

func (r *CategoryRepository) Create(category model.Category) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction")
		return 0, err
	}

	defer tx.Rollback()

	var categoryID int
	err = tx.QueryRow("INSERT INTO categories DEFAULT VALUES RETURNING id").Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	for _, translation := range category.Translations {
		_, err := tx.Exec(
			`INSERT INTO cat_translate (cat_id, lang_id, name)
			VALUES ($1, $2, $3)`, categoryID, translation.LangID, translation.Name,
		)
		if err != nil {
			return 0, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return categoryID, nil
}

// func (r *CategoryRepository) Update(id int, category model.Category) error {
// 	query := fmt.Sprintf("UPDATE %v SET name=$1 WHERE id=$2", CATEGORIES)
// 	_, err := r.DB.Exec(query, category.Name, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *CategoryRepository) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id = $1", CATEGORIES)
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// func (r *CategoryRepository) GetAll() ([]model.CategoryResponse, error) {
// 	var categories []model.CategoryResponse
// 	query := `SELECT c.id, c.name AS name, l.name AS language_name FROM categories AS c INNER JOIN languages AS l ON c.language_id = l.id`
// 	err := r.DB.Select(&categories, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return categories, nil
// }
