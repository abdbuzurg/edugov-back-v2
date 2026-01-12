package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeeWorkExperience handles creation requests for work experience.
func (h *Handlers) CreateEmployeeWorkExperience(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeWorkExperienceRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeWorkExperience(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeWorkExperiences lists work experiences for an employee.
func (h *Handlers) GetEmployeeWorkExperiences(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeWorkExperiencesByEmployeeIDAndLanguageCode(
		r.Context(),
		employeeID,
		lang,
	)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeWorkExperienceResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeWorkExperience handles update requests for work experience.
func (h *Handlers) UpdateEmployeeWorkExperience(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeWorkExperienceRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeWorkExperience(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeWorkExperience handles deletion requests for work experience.
func (h *Handlers) DeleteEmployeeWorkExperience(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeWorkExperience(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
