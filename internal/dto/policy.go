package dto

type PolicyRequest struct {
	Subject string `json:"subject" validate:"required"`
	Object  string `json:"object" validate:"required"`
	Action  string `json:"action" validate:"required"`
}

type PolicyGroupRequest struct {
	Subject string `json:"subject" validate:"required"`
	Object  string `json:"object" validate:"required"`
}

type PolicyResponse struct {
	Subject string `json:"subject"`
	Object  string `json:"object"`
	Action  string `json:"action"`
}
