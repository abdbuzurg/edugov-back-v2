package dto

import "time"

// CreateEmployeePublicationRequest represents input for creating a publication.
type CreateEmployeePublicationRequest struct {
	EmployeeID        int64  `json:"employeeID" validate:"required,min=1"`
	LanguageCode      string `json:"-" validate:"required,len=2"`
	PublicationTitle  string `json:"publicationTitle" validate:"required"`
	LinkToPublication string `json:"linkToPublication" validate:"required"`
}

// UpdateEmployeePublicationRequest represents input for updating a publication.
type UpdateEmployeePublicationRequest struct {
	ID                int64   `json:"id" validate:"required,min=1"`
	PublicationTitle  *string `json:"publicationTitle" validate:"omitempty"`
	LinkToPublication *string `json:"linkToPublication" validate:"omitempty"`
}

// EmployeePublicationResponse represents a publication in responses.
type EmployeePublicationResponse struct {
	ID                int64     `json:"id"`
	PublicationTitle  string    `json:"publicationTitle"`
	LinkToPublication string    `json:"linkToPublication"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
