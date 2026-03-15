package entity

import "github.com/google/uuid"

type Activity struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ActionType string    `json:"action_type"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	TargetName string    `json:"target_name"`
	CreatedAt  string    `json:"created_at"`
}
