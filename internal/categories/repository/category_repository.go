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

func (r *CategoryRepository) Create(category model.CategoryReq) (int, error) {
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

func (r *CategoryRepository) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id = $1", CATEGORIES)
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) GetAllByLangID(langId int) ([]model.CategoryRes, error) {
	var category []model.CategoryRes
	query := `
		SELECT 
			c.id, ct.name 
		FROM 
			categories AS c 
		INNER JOIN 
			cat_translate AS ct ON c.id=ct.cat_id 
		INNER JOIN 
			languages AS l ON l.id=ct.lang_id 
		WHERE lang_id = $1	
		`
	rows, err := r.DB.Query(query, langId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var res model.CategoryRes
		err := rows.Scan(&res.ID, &res.Name)
		if err != nil {
			return nil, err
		}
		category = append(category, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return category, nil
}
