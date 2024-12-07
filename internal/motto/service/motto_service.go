package service

import (
	"arassachylyk/internal/motto/model"
	"arassachylyk/internal/motto/repository"
)

type MottoService struct {
	repo *repository.MottoRepository
}

func NewYearService(repo *repository.MottoRepository) *MottoService {
	return &MottoService{repo: repo}
}

func (s *MottoService) Create(motto model.Motto) (int, error) {
	return s.repo.Create(motto)
}

func (s *MottoService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *MottoService) GetByID(id int) (model.Motto, error) {
	return s.repo.GetByID(id)
}
