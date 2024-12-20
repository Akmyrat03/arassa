package service

import (
	"arassachylyk/internal/categories/model"
	"arassachylyk/internal/categories/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(category model.CategoryReq) (int, error) {
	return s.repo.Create(category)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) GetAllByLangID(langId int) ([]model.CategoryRes, error) {
	return s.repo.GetAllByLangID(langId)
}
