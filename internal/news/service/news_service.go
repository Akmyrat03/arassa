package service

import (
	"arassachylyk/internal/news/model"
	"arassachylyk/internal/news/repository"
)

type NewsService struct {
	repo *repository.NewsRepository
}

func NewService(repo *repository.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) Create(news model.News) (int, error) {
	return s.repo.Create(news)
}

func (s *NewsService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *NewsService) GetNewsByID(id int) (model.News, error) {
	return s.repo.GetByID(id)
}

func (s *NewsService) GetAllNewsByLangID(langID, page, limit int) ([]model.NewsLang, error) {
	return s.repo.GetAllNewsByLangID(langID, page, limit)
}

func (s *NewsService) GetAllNewsByLangAndCategory(langID, categoryID, limit, page int) ([]model.NewsLang, error) {
	return s.repo.GetAllNewsByLangAndCategory(langID, categoryID, page, limit)
}
