package handlers

import (
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

func (h *Handlers) CreateEmployeePatent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeePatentRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	req.LanguageCode = middleware.MustLang(r.Context())

	resp, err := h.service.CreateEmployeePatent(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handlers) GetEmployeePatents(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeePatentsByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeePatentResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) UpdateEmployeePatent(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeePatentRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.UpdateEmployeePatent(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) DeleteEmployeePatent(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeePatent(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
