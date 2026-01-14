package dto

import "time"

// CreateEmployeePublicationRequest represents input for creating a publication.
type CreateEmployeePublicationRequest struct {
	EmployeeID          int64   `json:"employeeID" validate:"required,min=1"`
	LanguageCode        string  `json:"-" validate:"required,len=2"`
	RfPublicationTypeID int64   `json:"rfPublicationTypeID" validate:"required,min=1"`
	Name                string  `json:"name" validate:"required"`
	Type                string  `json:"type" validate:"required"`
	Authors             *string `json:"authors" validate:"omitempty"`
	JournalName         *string `json:"journalName" validate:"omitempty"`
	Volume              *string `json:"volume" validate:"omitempty"`
	Number              *string `json:"number" validate:"omitempty"`
	Pages               *string `json:"pages" validate:"omitempty"`
	Year                *int32  `json:"year" validate:"omitempty,min=0"`
	Link                string  `json:"link" validate:"required"`
}

// UpdateEmployeePublicationRequest represents input for updating a publication.
type UpdateEmployeePublicationRequest struct {
	ID                  int64   `json:"id" validate:"required,min=1"`
	RfPublicationTypeID *int64  `json:"rfPublicationTypeID" validate:"omitempty,min=1"`
	Name                *string `json:"name" validate:"omitempty"`
	Type                *string `json:"type" validate:"omitempty"`
	Authors             *string `json:"authors" validate:"omitempty"`
	JournalName         *string `json:"journalName" validate:"omitempty"`
	Volume              *string `json:"volume" validate:"omitempty"`
	Number              *string `json:"number" validate:"omitempty"`
	Pages               *string `json:"pages" validate:"omitempty"`
	Year                *int32  `json:"year" validate:"omitempty,min=0"`
	Link                *string `json:"link" validate:"omitempty"`
}

// EmployeePublicationResponse represents a publication in responses.
type EmployeePublicationResponse struct {
	ID                  int64     `json:"id"`
	RfPublicationTypeID int64     `json:"rfPublicationTypeID"`
	Name                string    `json:"name"`
	Type                string    `json:"type"`
	Authors             *string   `json:"authors,omitempty"`
	JournalName         *string   `json:"journalName,omitempty"`
	Volume              *string   `json:"volume,omitempty"`
	Number              *string   `json:"number,omitempty"`
	Pages               *string   `json:"pages,omitempty"`
	Year                *int32    `json:"year,omitempty"`
	Link                string    `json:"link"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}
