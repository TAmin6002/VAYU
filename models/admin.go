package models

import "time"

type Admin struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PublicAdmin is the shape returned to clients (never includes the hash).
type PublicAdmin struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (a Admin) ToPublic() PublicAdmin {
	return PublicAdmin{
		ID:       a.ID,
		FullName: a.FullName,
		Email:    a.Email,
	}
}
