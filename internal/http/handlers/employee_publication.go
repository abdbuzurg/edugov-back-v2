package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeePublication handles creation requests for publications.
func (h *Handlers) CreateEmployeePublication(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeePublicationRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeePublication(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeePublications lists publications for an employee.
func (h *Handlers) GetEmployeePublications(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeePublicationByEmployeeIDAndLanguageCode(
		r.Context(),
		employeeID,
		lang,
	)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeePublicationResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeePublication handles update requests for publications.
func (h *Handlers) UpdateEmployeePublication(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeePublicationRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeePublication(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeePublication handles deletion requests for publications.
func (h *Handlers) DeleteEmployeePublication(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeePublication(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
