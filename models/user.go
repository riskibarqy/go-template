package models

// User models
type User struct {
	ID             int     `json:"id" db:"id"`
	Name           string  `json:"name" db:"name"`
	Email          string  `json:"email" db:"email"`
	Password       string  `json:"password,omitempty" db:"password"`
	Token          *string `json:"token,omitempty" db:"token"`
	TokenExpiredAt *int    `json:"tokenExpiredAt,omitempty" db:"token_expired_at"`
	CreatedAt      int     `json:"createdAt" db:"created_at"`
	UpdatedAt      *int    `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt      *int    `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (u *User) ForPublic() {
	u.Password = ""
	u.Token = nil
	u.TokenExpiredAt = nil
	u.UpdatedAt = nil
	u.DeletedAt = nil
}
