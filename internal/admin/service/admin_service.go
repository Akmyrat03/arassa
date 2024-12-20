package service

import (
	"arassachylyk/internal/admin/model"
	"arassachylyk/internal/admin/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	salt       = "sdnksnk34welfw23"
	signingKey = []byte("###%4544566656")
)

type tokenClaims struct {
	jwt.StandardClaims

	Admin_id int    `json:"admin_id"`
	Username string `json:"username"`
}

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) CreateAdmin(admin *model.Admin) (int, error) {
	admin.Password = GeneratePasswordHash(admin.Password)
	return s.repo.Create(admin)
}

func (s *AdminService) GetAdmin(username, password string) (model.Admin, error) {
	return s.repo.GetAdmin(username, password)
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()

	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AdminService) GenerateToken(username, password string) (string, error) {
	admin, err := s.repo.GetAdmin(username, GeneratePasswordHash(password))
	if err != nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		admin.ID,
		admin.Username,
	})

	return token.SignedString(signingKey)
}

// Validate token
func (s *AdminService) ValidateToken(tokenString string) (string, error) {
	claims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Username, nil
}
