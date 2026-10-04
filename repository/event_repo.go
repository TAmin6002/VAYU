package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"university-events/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, e *models.Event) error {
	query := `
		INSERT INTO events (title, description, event_time, location, capacity, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		e.Title, e.Description, e.EventTime, e.Location, e.Capacity, e.CreatedBy,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}


var ErrCapacityBelowRegistered = errors.New("capacity is below the current registered count")

func (r *EventRepository) Update(ctx context.Context, id string, in models.UpdateEventInput) (*models.Event, error) {
	registeredCount, err := r.registeredCount(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Capacity < registeredCount {
		return nil, ErrCapacityBelowRegistered
	}

	query := `
		UPDATE events
		SET title = $1, description = $2, event_time = $3, location = $4, capacity = $5, updated_at = now()
		WHERE id = $6
		RETURNING id, title, description, event_time, location, capacity, created_by, created_at, updated_at`

	var e models.Event
	err = r.db.QueryRowContext(ctx, query,
		in.Title, in.Description, in.EventTime, in.Location, in.Capacity, id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.EventTime, &e.Location, &e.Capacity, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("update event: %w", err)
	}
	return &e, nil
}


func (r *EventRepository) registeredCount(ctx context.Context, eventID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM participations WHERE event_id = $1 AND status = 'registered'`,
		eventID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count registered participants: %w", err)
	}
	return count, nil
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
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

func (r *EventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	query := `
		SELECT id, title, description, event_time, location, capacity, created_by, created_at, updated_at
		FROM events WHERE id = $1`

	var e models.Event
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventTime, &e.Location, &e.Capacity, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get event by id: %w", err)
	}
	return &e, nil
}


func (r *EventRepository) List(ctx context.Context) ([]models.EventWithStats, error) {
	query := `
		SELECT
			e.id, e.title, e.description, e.event_time, e.location, e.capacity,
			e.created_by, e.created_at, e.updated_at,
			COALESCE(p.registered_count, 0) AS registered_count
		FROM events e
		LEFT JOIN (
			SELECT event_id, COUNT(*) AS registered_count
			FROM participations
			WHERE status = 'registered'
			GROUP BY event_id
		) p ON p.event_id = e.id
		ORDER BY e.event_time ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	var events []models.EventWithStats
	for rows.Next() {
		var ev models.EventWithStats
		if err := rows.Scan(
			&ev.ID, &ev.Title, &ev.Description, &ev.EventTime, &ev.Location, &ev.Capacity,
			&ev.CreatedBy, &ev.CreatedAt, &ev.UpdatedAt, &ev.RegisteredCount,
		); err != nil {
			return nil, fmt.Errorf("scan event row: %w", err)
		}
		ev.IsFull = ev.RegisteredCount >= ev.Capacity
		events = append(events, ev)
	}
	return events, rows.Err()
}

