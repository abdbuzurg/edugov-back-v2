package handlers

import (
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
)

func (h *Handlers) GetEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	resp, err := h.service.GetEmployeeDetailsByEmployeeID(r.Context(), employeeID)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) UpdateEmployeeDetails(w http.ResponseWriter, r *http.Request) {
	var body dto.UpdateFullEmployeeData
	if err := DecodeJSON(r, &body); err != nil {
		WriteError(w, r, err)
		return
	}

	if len(body.Data) == 0 {
		WriteError(w, r, apperr.Validation("validation failed", map[string]string{
			"data": "required",
		}))
		return
	}

	lang := middleware.MustLang(r.Context())

	// Enforce employeeID + languageCode from context (FE controls Accept-Language)
	// so client can't mix employee IDs or languages in one request.
	for i := range body.Data {
		body.Data[i].LanguageCode = lang
	}

	resp, err := h.service.UpdateEmployeeDetails(r.Context(), body.Data)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}
