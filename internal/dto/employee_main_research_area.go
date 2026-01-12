package dto

import "time"

// CreateEmployeeMainResearchAreaRequest represents input for creating a main research area.
type CreateEmployeeMainResearchAreaRequest struct {
	EmployeeID   int64                                `json:"employeeID" validate:"required,min=1"`
	LanguageCode string                               `json:"languageCode" validate:"required,len=2"`
	Area         string                               `json:"area" validate:"required"`
	Discipline   string                               `json:"discipline" validate:"required"`
	KeyTopics    []*CreateResearchAreaKeyTopicRequest `json:"keyTopics" validate:"required,dive"`
}

// UpdateEmployeeMainResearchAreaRequest represents input for updating a main research area.
type UpdateEmployeeMainResearchAreaRequest struct {
	ID         int64                                `json:"id" validate:"required,min=1"`
	Discipline *string                              `json:"discipline" validate:"omitempty"`
	Area       *string                              `json:"area" validate:"omitempty"`
	KeyTopics  []*UpdateResearchAreaKeyTopicRequest `json:"keyTopics" required:"omitempty,dive"`
}

// CreateResearchAreaKeyTopicRequest represents input for a key topic.
type CreateResearchAreaKeyTopicRequest struct {
	LanguageCode  string `json:"languageCode" validate:"required,len=2"`
	KeyTopicTitle string `json:"keyTopicTitle" validate:"required"`
}

// UpdateResearchAreaKeyTopicRequest represents input for updating a key topic.
type UpdateResearchAreaKeyTopicRequest struct {
	ID            int64   `json:"id" validate:"required,min=1"`
	KeyTopicTitle *string `json:"keyTopicTitle" validate:"omitempty"`
}

// EmployeeMainResearchAreaResponse represents a research area with key topics.
type EmployeeMainResearchAreaResponse struct {
	ID         int64                           `json:"id"`
	Discipline string                          `json:"discipline"`
	Area       string                          `json:"area"`
	KeyTopics  []*ResearchAreaKeyTopicResponse `json:"keyTopics,omitempty"`
	CreatedAt  time.Time                       `json:"createdAt"`
	UpdatedAt  time.Time                       `json:"updatedAt"`
}

// ResearchAreaKeyTopicResponse represents a key topic in responses.
type ResearchAreaKeyTopicResponse struct {
	ID            int64     `json:"id"`
	KeyTopicTitle string    `json:"keyTopicTitle"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
