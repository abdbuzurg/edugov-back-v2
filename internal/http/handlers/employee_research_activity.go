package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeeResearchActivity handles creation requests for research activities.
func (h *Handlers) CreateEmployeeResearchActivity(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeResearchActivityRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeResearchActivity(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeResearchActivities lists research activities for an employee.
func (h *Handlers) GetEmployeeResearchActivities(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeResearchActivitiesByEmployeeIDAndLanguageCode(
		r.Context(),
		employeeID,
		lang,
	)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeResearchActivityResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeResearchActivity handles update requests for research activities.
func (h *Handlers) UpdateEmployeeResearchActivity(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeResearchActivityRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeResearchActivity(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeResearchActivity handles deletion requests for research activities.
func (h *Handlers) DeleteEmployeeResearchActivity(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeResearchActivity(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
