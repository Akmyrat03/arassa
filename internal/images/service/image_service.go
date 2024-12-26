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

func (s *ImageService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ImageService) GetImageByTitleID(id int) ([]string, error) {
	return s.repo.GetImagePathsByTitleID(id)
}

func (s *ImageService) GetAll(langID int) ([]model.Image, error) {
	return s.repo.GetAllImages(langID)
}

func (s *ImageService) GetPaginatedImg(langID, page, limit int) ([]model.Image, error) {
	return s.repo.GetPaginatedImages(langID, page, limit)
}
