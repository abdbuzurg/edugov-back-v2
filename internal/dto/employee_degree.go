package dto

import "time"

// CreateEmployeeDegreeRequest represents input for creating a degree.
type CreateEmployeeDegreeRequest struct {
	EmployeeID         int64     `json:"employeeID" validate:"required,min=1"`
	LanguageCode       string    `json:"-" validate:"required,len=2"`
	RfInstitutionID    int64     `json:"rfInstitutionID" validate:"required,min=1"`
	DegreeLevel        string    `json:"degreeLevel" validate:"required"`
	InstitutionName    string    `json:"institutionName" validate:"required"`
	Speciality         string    `json:"speciality" validate:"required"`
	DateStart          time.Time `json:"dateStart" validate:"required"`
	DateEnd            time.Time `json:"dateEnd" validate:"required"`
	GivenBy            string    `json:"givenBy" validate:"required"`
	DateDegreeRecieved time.Time `json:"dateDegreeRecieved" validate:"required"`
}

// UpdateEmployeeDegreeRequest represents input for updating a degree.
type UpdateEmployeeDegreeRequest struct {
	ID                 int64      `json:"id" validate:"required,min=1"`
	RfInstitutionID    *int64     `json:"rfInstitutionID" validate:"omitempty,min=1"`
	DegreeLevel        *string    `json:"degreeLevel" validate:"omitempty"`
	InstitutionName    *string    `json:"institutionName" validate:"omitempty"`
	Speciality         *string    `json:"speciality" validate:"omitempty"`
	DateStart          *time.Time `json:"dateStart" validate:"omitempty"`
	DateEnd            *time.Time `json:"dateEnd" validate:"omitempty"`
	GivenBy            *string    `json:"givenBy" validate:"omitempty"`
	DateDegreeRecieved *time.Time `json:"dateDegreeRecieved" validate:"omitempty"`
}

// EmployeeDegreeResponse represents a degree record in responses.
type EmployeeDegreeResponse struct {
	ID                 int64     `json:"id"`
	RfInstitutionID    int64     `json:"rfInstitutionID"`
	DegreeLevel        string    `json:"degreeLevel"`
	InstitutionName    string    `json:"institutionName"`
	Speciality         string    `json:"speciality"`
	DateStart          time.Time `json:"dateStart"`
	DateEnd            time.Time `json:"dateEnd"`
	GivenBy            string    `json:"givenBy"`
	DateDegreeRecieved time.Time `json:"dateDegreeRecieved"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
