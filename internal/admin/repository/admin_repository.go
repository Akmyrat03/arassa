package repository

import (
	"arassachylyk/internal/admin/model"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	Admin = "admin"
)

type AdminRepository struct {
	DB *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{DB: db}
}

func (r *AdminRepository) Create(admin *model.Admin) (int, error) {
	var id int
	err := r.DB.QueryRow(signUpQuery, admin.Username, admin.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AdminRepository) GetAdmin(username, password string) (model.Admin, error) {
	var user model.Admin
	err := r.DB.Get(&user, getAdminQuery, username, password)
	if err != nil {
		return model.Admin{}, errors.New("incorrect username or password")
	}

	return user, nil
}

func (r *AdminRepository) GetUserByField(field, value string) (model.Admin, error) {
	if field != "username" {
		return model.Admin{}, fmt.Errorf("unsupported field: %s", field)
	}

	query := fmt.Sprintf("SELECT id, username, password FROM %v WHERE %s= $1", Admin, field)
	var user model.Admin
	err := r.DB.Get(&user, query, value)
	if err != nil {
		return model.Admin{}, err
	}

	return user, nil
}
