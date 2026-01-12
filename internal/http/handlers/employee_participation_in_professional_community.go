package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

// CreateEmployeeParticipationInProfessionalCommunity handles creation requests for communities.
func (h *Handlers) CreateEmployeeParticipationInProfessionalCommunity(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeParticipationInProfessionalCommunityRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeParticipationInProfessionalCommunity(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeParticipationInProfessionalCommunities lists community participation for an employee.
func (h *Handlers) GetEmployeeParticipationInProfessionalCommunities(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeParticipationInProfessionalCommunitiesByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeParticipationInProfessionalCommunityResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeParticipationInProfessionalCommunity handles update requests for communities.
func (h *Handlers) UpdateEmployeeParticipationInProfessionalCommunity(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeParticipationInProfessionalCommunityRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeParticipationInProfessionalCommunity(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeParticipationInProfessionalCommunity handles deletion requests for communities.
func (h *Handlers) DeleteEmployeeParticipationInProfessionalCommunity(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeParticipationInProfessionalCommunity(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
