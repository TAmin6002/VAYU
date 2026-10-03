package models

import "time"

// Student represents a row in the students table.
type Student struct {
	ID            string    `json:"id"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	StudentNumber string    `json:"student_number"`
	PasswordHash  string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// PublicStudent is the shape returned to clients (never includes the hash).
type PublicStudent struct {
	ID            string `json:"id"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	StudentNumber string `json:"student_number"`
}

func (s Student) ToPublic() PublicStudent {
	return PublicStudent{
		ID:            s.ID,
		FullName:      s.FullName,
		Email:         s.Email,
		StudentNumber: s.StudentNumber,
	}
}

// SignupInput is the payload a student submits to create an account.
type SignupInput struct {
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	StudentNumber string `json:"student_number"`
	Password      string `json:"password"`
}

// MyRegistration is one active registration of the logged-in student.
type MyRegistration struct {
	EventID      string    `json:"event_id"`
	RegisteredAt time.Time `json:"registered_at"`
}
