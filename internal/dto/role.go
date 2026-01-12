package dto

import "time"

type RoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type RoleResponse struct {
	ID        string    `josn:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
