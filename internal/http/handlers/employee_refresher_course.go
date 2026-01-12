package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeeRefresherCourse handles creation requests for refresher courses.
func (h *Handlers) CreateEmployeeRefresherCourse(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeRefresherCourseRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeRefresherCourse(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeRefresherCourses lists refresher courses for an employee.
func (h *Handlers) GetEmployeeRefresherCourses(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeRefresherCoursesByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeRefresherCourseResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeRefresherCourse handles update requests for refresher courses.
func (h *Handlers) UpdateEmployeeRefresherCourse(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeRefresherCourseRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeRefresherCourse(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeRefresherCourse handles deletion requests for refresher courses.
func (h *Handlers) DeleteEmployeeRefresherCourse(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeRefresherCourse(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
