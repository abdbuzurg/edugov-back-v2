package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeeParticipationInEvent handles creation requests for events.
func (h *Handlers) CreateEmployeeParticipationInEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeParticipationInEventRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeParticipationInEvent(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeParticipationInEvents lists event participation for an employee.
func (h *Handlers) GetEmployeeParticipationInEvents(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeParticipationInEventsByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeParticipationInEventResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeParticipationInEvent handles update requests for events.
func (h *Handlers) UpdateEmployeeParticipationInEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeParticipationInEventRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeParticipationInEvent(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeParticipationInEvent handles deletion requests for events.
func (h *Handlers) DeleteEmployeeParticipationInEvent(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeParticipationInEvent(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
