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
	return s.repo.GetNewsByID(id)
}

func (s *NewsService) GetAll() ([]model.News, error) {
	return s.repo.GetAll()
}

func (s *NewsService) GetByCategoryID(id int) (model.News, error) {
	return s.repo.GetByCategoryID(id)
}
