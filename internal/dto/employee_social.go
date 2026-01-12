package dto

import "time"

// CreateEmployeeSocialRequest represents input for creating a social link.
type CreateEmployeeSocialRequest struct {
	EmployeeID   int64  `json:"employeeID" validate:"required,min=1"`
	SocialName   string `json:"socialName" validate:"required"`
	LinkToSocial string `json:"linkToSocial" validate:"required"`
}

// UpdateEmployeeSocialRequest represents input for updating a social link.
type UpdateEmployeeSocialRequest struct {
	ID           int64   `json:"id" validate:"required,min=1"`
	SocialName   *string `json:"socialName" validate:"omitempty"`
	LinkToSocial *string `json:"linkToSocial" validate:"omitempty"`
}

// EmployeeSocialResponse represents a social link in responses.
type EmployeeSocialResponse struct {
	ID           int64     `json:"id"`
	SocialName   string    `json:"socialName"`
	LinkToSocial string    `json:"linkToSocial"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
