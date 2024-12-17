package service

import (
	"arassachylyk/internal/images/model"
	"arassachylyk/internal/images/repository"
)

type ImageService struct {
	repo *repository.ImageRepository
}

func NewImageService(repo *repository.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) Create(title model.Title) (int, error) {
	return s.repo.Create(title)
}
