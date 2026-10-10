package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"university-events/models"
)

type ParticipationRepository struct {
	db *sql.DB
}

func NewParticipationRepository(db *sql.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

func (r *ParticipationRepository) Register(ctx context.Context, studentID, fullName, email, studentNumber, eventID string) (*models.Participation, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // no-op if committed


	var capacity int
	err = tx.QueryRowContext(ctx,
		`SELECT capacity FROM events WHERE id = $1 FOR UPDATE`, eventID,
	).Scan(&capacity)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("lock event: %w", err)
	}

	var registeredCount int
	err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM participations WHERE event_id = $1 AND status = 'registered'`, eventID,
	).Scan(&registeredCount)
	if err != nil {
		return nil, fmt.Errorf("count participations: %w", err)
	}

	if registeredCount >= capacity {
		return nil, models.ErrEventFull
	}

	var p models.Participation
	insertQuery := `
		INSERT INTO participations (full_name, email, student_number, event_id, student_id, status)
		VALUES ($1, $2, $3, $4, $5, 'registered')
		RETURNING id, event_id, full_name, email, student_number, status, registered_at`

	err = tx.QueryRowContext(ctx, insertQuery, fullName, email, studentNumber, eventID, studentID).
		Scan(&p.ID, &p.EventID, &p.FullName, &p.Email, &p.StudentNumber, &p.Status, &p.RegisteredAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return nil, models.ErrAlreadyRegistered
		}
		return nil, fmt.Errorf("insert participation: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &p, nil
}


func (r *ParticipationRepository) Cancel(ctx context.Context, eventID, studentID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE participations SET status = 'cancelled' WHERE event_id = $1 AND student_id = $2 AND status = 'registered'`,
		eventID, studentID,
	)
	if err != nil {
		return fmt.Errorf("cancel participation: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ParticipationRepository) ListByStudent(ctx context.Context, studentID string) ([]models.MyRegistration, error) {
	query := `
		SELECT p.event_id, p.registered_at
		FROM participations p
		JOIN events e ON e.id = p.event_id
		WHERE p.student_id = $1 AND p.status = 'registered'
		ORDER BY e.event_time ASC`

	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("list student registrations: %w", err)
	}
	defer rows.Close()

	regs := []models.MyRegistration{}
	for rows.Next() {
		var m models.MyRegistration
		if err := rows.Scan(&m.EventID, &m.RegisteredAt); err != nil {
			return nil, fmt.Errorf("scan registration row: %w", err)
		}
		regs = append(regs, m)
	}
	return regs, rows.Err()
}

// ListParticipants returns everyone registered for an event (for the admin view).
func (r *ParticipationRepository) ListParticipants(ctx context.Context, eventID string) ([]models.Participant, error) {
	query := `
		SELECT full_name, email, student_number, status, registered_at
		FROM participations
		WHERE event_id = $1
		ORDER BY registered_at ASC`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("list participants: %w", err)
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.FullName, &p.Email, &p.StudentNumber, &p.Status, &p.RegisteredAt); err != nil {
			return nil, fmt.Errorf("scan participant row: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}
