package dto

import "time"

type ListUsersInput struct {
	PagingInput
	Username string `form:"username"`
}

type ListUsersOutput struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
