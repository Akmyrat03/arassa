package repository

import (
	"arassachylyk/internal/news/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	NEWS = "news"
)

type NewsRepository struct {
	DB *sqlx.DB
}

func NewRepository(DB *sqlx.DB) *NewsRepository {
	return &NewsRepository{DB: DB}
}

// func (r *NewsRepository) Create(news model.News) (int, error) {
// 	var id int
// 	query := fmt.Sprintf("INSERT INTO %v (category_id, title, description, image, created_at) VALUES ($1, $2, $3, $4, Now()) RETURNING id", NEWS)
// 	rows := r.DB.QueryRow(query, news.CategoryID, news.Title, news.Description, news.ImageURL)
// 	if err := rows.Scan(&id); err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

func (r *NewsRepository) Create(news model.News) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Println("Error starting transaction")
		return 0, err
	}

	defer tx.Rollback()

	var newsId int
	query := `INSERT INTO news (image, category_id) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, news.ImageURL, news.CategoryID).Scan(&newsId)
	if err != nil {
		return 0, err
	}

	for _, translation := range news.Translations {
		query := `INSERT INTO news_translate (news_id, lang_id, title, description) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(query, newsId, translation.LangID, translation.Title, translation.Description)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return newsId, nil
}

func (r *NewsRepository) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id=$1", NEWS)
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *NewsRepository) GetNewsByID(id int) (model.News, error) {
	var news model.News
	query := `SELECT n.id, n.title, n.description, n.image, n.created_at, c.name AS category_name FROM news AS n INNER JOIN categories AS c ON n.category_id = c.id WHERE n.id=$1`
	err := r.DB.Get(&news, query, id)
	if err != nil {
		return model.News{}, err
	}

	return news, nil
}

func (r *NewsRepository) GetAll() ([]model.News, error) {
	var news []model.News
	query := `SELECT n.id, n.title, n.description, n.image, n.created_at, c.id AS category_id, c.name AS category_name FROM news AS n INNER JOIN categories AS c ON n.category_id = c.id`
	err := r.DB.Select(&news, query)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (r *NewsRepository) GetByCategoryID(id int) (model.News, error) {
	var news model.News
	query := `SELECT n.id, n.title, n.description, n.image, n.created_at, c.id AS category_id, c.name AS category_name FROM news AS n INNER JOIN categories AS c ON n.category_id = c.id WHERE c.id=$1`
	err := r.DB.Get(&news, query, id)
	if err != nil {
		return model.News{}, err
	}
	return news, nil
}
