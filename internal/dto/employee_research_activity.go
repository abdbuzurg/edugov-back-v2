package dto

import "time"

// CreateEmployeeResearchActivityRequest represents input for creating research activity.
type CreateEmployeeResearchActivityRequest struct {
	EmployeeID            int64  `json:"employeeID" validate:"required,min=1"`
	LanguageCode          string `json:"-" validate:"required,len=2"`
	ResearchActivityTitle string `json:"researchActivityTitle" validate:"required"`
	EmployeeRole          string `json:"employeeRole" validate:"required"`
}

// UpdateEmployeeResearchActivityRequest represents input for updating research activity.
type UpdateEmployeeResearchActivityRequest struct {
	ID                    int64   `json:"id" validate:"required,min=1"`
	ResearchActivityTitle *string `json:"researchActivityTitle" validate:"omitempty"`
	EmployeeRole          *string `json:"employeeRole" validate:"omitempty"`
}

// EmployeeResearchActivityResponse represents research activity in responses.
type EmployeeResearchActivityResponse struct {
	ID                    int64     `json:"id"`
	ResearchActivityTitle string    `json:"researchActivityTitle"`
	EmployeeRole          string    `json:"employeeRole"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
