package models

import "time"

type ParticipationStatus string

const (
	StatusRegistered ParticipationStatus = "registered"
	StatusCancelled  ParticipationStatus = "cancelled"
)


type Participation struct {
	ID            string              `json:"id"`
	EventID       string              `json:"event_id"`
	FullName      string              `json:"full_name"`
	Email         string              `json:"email"`
	StudentNumber string              `json:"student_number"`
	Status        ParticipationStatus `json:"status"`
	RegisteredAt  time.Time           `json:"registered_at"`
}

// RegisterInput is the payload a visitor submits on the event detail page.
type RegisterInput struct {
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	StudentNumber string `json:"student_number"`
}


type Participant struct {
	FullName      string              `json:"full_name"`
	Email         string              `json:"email"`
	StudentNumber string              `json:"student_number"`
	Status        ParticipationStatus `json:"status"`
	RegisteredAt  time.Time           `json:"registered_at"`
}


var ErrEventFull = &DomainError{Code: "event_full", Message: "capacity reached for this event"}

// ErrAlreadyRegistered is returned when a user tries to register twice.
var ErrAlreadyRegistered = &DomainError{Code: "already_registered", Message: "user is already registered for this event"}

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}
