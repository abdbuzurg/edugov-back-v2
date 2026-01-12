package handlers

import (
	"edugov-back-v2/internal/dto"
	"net/http"
)

// CreateEmployeeSocial handles creation requests for social links.
func (h *Handlers) CreateEmployeeSocial(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeSocialRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.CreateEmployeeSocial(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetEmployeeSocials lists social links for an employee.
func (h *Handlers) GetEmployeeSocials(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.GetEmployeeSocialByEmployeeID(r.Context(), employeeID)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeSocialResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

// UpdateEmployeeSocial handles update requests for social links.
func (h *Handlers) UpdateEmployeeSocial(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeSocialRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeSocial(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

// DeleteEmployeeSocial handles deletion requests for social links.
func (h *Handlers) DeleteEmployeeSocial(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeSocial(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
