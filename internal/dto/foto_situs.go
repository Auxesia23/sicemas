package dto

import "github.com/google/uuid"

// Taruh ini di folder entity atau dto kamu
type FotoSitusPayload struct {
	SitusID  uuid.UUID `db:"situs_id"`
	ImageURL string    `db:"image_url"`
	PublicID string    `db:"public_id"`
}

type UploadResult struct {
	SecureURL string
	PublicID  string
}

type FotoResponse struct {
	ID       uuid.UUID `json:"id" db:"id"`
	ImageURL string    `json:"image_url" db:"image_url"`
	PublicID string    `json:"public_id" db:"public_id"`
}
