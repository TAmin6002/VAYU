package models

import "time"

// Admin represents a row in the admins table. There is exactly one
// account type in this system — there are no visitor/user accounts and
// no roles. The single admin account is created manually in the
// database (see cmd/hashpassword and the README); regular visitors
// never get an account at all, they just register for events directly
// (see models.Participation).
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
