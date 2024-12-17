package repository

import (
	"arassachylyk/internal/contact/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type ContactRepository struct {
	DB *sqlx.DB
}

func NewContactRepository(db *sqlx.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

func (r *ContactRepository) SaveMessage(ctx context.Context, message model.ContactMessage) error {
	query := `INSERT INTO contact_messages (name, email, message, phone_number) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, message.Name, message.Email, message.Message, message.PhoneNumber)
	return err
}
