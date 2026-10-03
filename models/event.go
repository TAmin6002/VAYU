package models

import "time"

// Event represents a row in the events table.
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventWithStats extends Event with participation numbers, used in list views.
type EventWithStats struct {
	Event
	RegisteredCount int  `json:"registered_count"`
	IsFull          bool `json:"is_full"`
}

// CreateEventInput is the payload accepted from a group admin creating an event.
type CreateEventInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
}

// UpdateEventInput is the payload accepted when editing an event.
type UpdateEventInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
}
