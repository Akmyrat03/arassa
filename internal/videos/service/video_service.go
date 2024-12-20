package service

import (
	"arassachylyk/internal/videos/model"
	"arassachylyk/internal/videos/repository"
)

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(repo *repository.VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) UploadVideos(title model.Title) (int, error) {
	return s.repo.Upload(title)
}

func (s *VideoService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *VideoService) GetVideoPaths(id int) ([]string, error) {
	return s.repo.GetVideoPathsByID(id)
}

func (s *VideoService) GetAll(id int) ([]model.Video, error) {
	return s.repo.GetAllVideos(id)
}
