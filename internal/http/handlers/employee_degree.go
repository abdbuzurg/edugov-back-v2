package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

func (h *Handlers) CreateEmployeeDegree(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeDegreeRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeeDegree(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handlers) UpdateEmployeeDegree(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeDegreeRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeeDegree(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) DeleteEmployeeDegree(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeDegree(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetEmployeeDegrees(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeDegreesByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeDegreeResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}
