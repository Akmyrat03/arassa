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

	q := `DELETE FROM news_translate WHERE id=$1`
	_, err = r.DB.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *NewsRepository) GetByID(id int) (model.News, error) {
	var news model.News
	translations := []model.Translation{}
	query := `
		SELECT n.id AS news_id, n.category_id, n.image, n.created_at, nt.title, nt.description, nt.lang_id 
		FROM news AS n 
		LEFT JOIN news_translate AS nt ON nt.news_id = n.id 
		WHERE n.id=$1
	`
	rows, err := r.DB.Query(query, id)
	if err != nil {
		return news, err
	}
	defer rows.Close()

	for rows.Next() {
		var translation model.Translation
		err := rows.Scan(
			&news.ID, &news.CategoryID, &news.ImageURL, &news.CreatedAt, &translation.Title, &translation.Description, &translation.LangID,
		)
		if err != nil {
			return news, err
		}
		translations = append(translations, translation)
	}
	news.Translations = translations
	return news, nil
}

func (r *NewsRepository) GetAllNewsByLangID(langID int) ([]model.NewsLang, error) {
	var newsListTKM []model.NewsLang
	query := `
	SELECT 
		n.id AS news_id, nt.title, ct.name AS category, nt.description, n.image, n.created_at 
	FROM 
		news_translate AS nt 
	INNER JOIN 
		news AS n ON nt.news_id= n.id 
	INNER JOIN 
		cat_translate AS ct ON n.category_id=ct.cat_id AND nt.lang_id=ct.lang_id 
	WHERE 
		nt.lang_id = $1
	`
	rows, err := r.DB.Query(query, langID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var news model.NewsLang
		err := rows.Scan(
			&news.ID, &news.Title, &news.CategoryName, &news.Description, &news.ImageURL, &news.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		newsListTKM = append(newsListTKM, news)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return newsListTKM, nil
}

func (r *NewsRepository) GetAllNewsByLangAndCategory(langID, categoryID int) ([]model.NewsLang, error) {
	var newsList []model.NewsLang
	query := `
	SELECT 
		n.id AS news_id, nt.title, ct.name AS category, nt.description, n.image, n.created_at 
	FROM 
		news_translate AS nt 
	INNER JOIN 
		news AS n ON nt.news_id = n.id 
	INNER JOIN 
		cat_translate AS ct ON n.category_id = ct.cat_id AND nt.lang_id = ct.lang_id 
	WHERE 
		nt.lang_id = $1 AND n.category_id = $2
	`
	rows, err := r.DB.Query(query, langID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var news model.NewsLang
		err := rows.Scan(
			&news.ID, &news.Title, &news.CategoryName, &news.Description, &news.ImageURL, &news.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}
