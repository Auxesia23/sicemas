package dto

import (
	"time"

	"github.com/google/uuid"
)

type ActivityResponse struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User      string    `json:"user" db:"actor_name"`
	Action    string    `json:"action" db:"action_type"`
	Target    string    `json:"target" db:"target_name"`
	Entity    string    `json:"entity" db:"entity_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
