package dto

import "time"

// CreateEmployeeRefresherCourseRequest represents input for creating a course.
type CreateEmployeeRefresherCourseRequest struct {
	EmployeeID   int64     `json:"employeeID" validate:"required,min=1"`
	LanguageCode string    `json:"-" validate:"required,len=2"`
	CourseTitle  string    `json:"courseTitle" validate:"required"`
	DateStart    time.Time `json:"dateStart" validate:"required"`
	DateEnd      time.Time `json:"dateEnd" validate:"required"`
}

// UpdateEmployeeRefresherCourseRequest represents input for updating a course.
type UpdateEmployeeRefresherCourseRequest struct {
	ID          int64      `json:"id" validate:"required,min=1"`
	CourseTitle *string    `json:"courseTitle" validate:"omitempty"`
	DateStart   *time.Time `json:"dateStart" validate:"omitempty"`
	DateEnd     *time.Time `json:"dateEnd" validate:"omitempty"`
}

// EmployeeRefresherCourseResponse represents a refresher course in responses.
type EmployeeRefresherCourseResponse struct {
	ID          int64     `json:"id"`
	CourseTitle string    `json:"courseTitle"`
	DateStart   time.Time `json:"dateStart"`
	DateEnd     time.Time `json:"dateEnd"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
