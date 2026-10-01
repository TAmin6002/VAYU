package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"university-events/models"
)

var ErrNotFound = errors.New("record not found")

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at, updated_at
		FROM admins WHERE email = $1`

	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&a.ID, &a.FullName, &a.Email, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get admin by email: %w", err)
	}
	return &a, nil
}

func (r *AdminRepository) GetByID(ctx context.Context, id string) (*models.Admin, error) {
	query := `
		SELECT id, full_name, email, password_hash, created_at, updated_at
		FROM admins WHERE id = $1`

	var a models.Admin
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.FullName, &a.Email, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get admin by id: %w", err)
	}
	return &a, nil
}
