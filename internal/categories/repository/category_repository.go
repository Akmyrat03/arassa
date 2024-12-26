package repository

import (
	"arassachylyk/internal/categories/model"
	"log"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository struct {
	DB *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Create(category model.CategoryReq) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction")
		return 0, err
	}

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && err == nil {
			log.Println("Error during transaction rollback:", rollbackErr)
		}
	}()

	var categoryID int

	err = tx.QueryRow(categoriesQuery).Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	for _, translation := range category.Translations {
		_, err := tx.Exec(
			categoryTranslateQuery, categoryID, translation.LangID, translation.Name,
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
	_, err := r.DB.Exec(deleteCategory, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) GetAllByLangID(langID int) ([]model.CategoryRes, error) {
	var category []model.CategoryRes

	rows, err := r.DB.Query(getAllCategoriesByLangID, langID)
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
