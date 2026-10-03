package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"university-events/models"
)

var (
	ErrEmailTaken         = errors.New("email is already in use")
	ErrStudentNumberTaken = errors.New("student number is already in use")
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Create inserts a new student account and fills in its generated fields.
func (r *StudentRepository) Create(ctx context.Context, s *models.Student) error {
	query := `
		INSERT INTO students (full_name, email, student_number, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		s.FullName, s.Email, s.StudentNumber, s.PasswordHash,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
		if strings.Contains(pqErr.Constraint, "student_number") {
			return ErrStudentNumberTaken
		}
		return ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

func (r *StudentRepository) GetByEmail(ctx context.Context, email string) (*models.Student, error) {
	return r.getOne(ctx, "email", email)
}

func (r *StudentRepository) GetByID(ctx context.Context, id string) (*models.Student, error) {
	return r.getOne(ctx, "id", id)
}

// getOne loads a single student by a column that is either "id" or "email"
// (never user-supplied, so building the query this way is safe).
func (r *StudentRepository) getOne(ctx context.Context, column, value string) (*models.Student, error) {
	query := `
		SELECT id, full_name, email, student_number, password_hash, created_at, updated_at
		FROM students WHERE ` + column + ` = $1`

	var s models.Student
	err := r.db.QueryRowContext(ctx, query, value).Scan(
		&s.ID, &s.FullName, &s.Email, &s.StudentNumber, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get student by %s: %w", column, err)
	}
	return &s, nil
}
