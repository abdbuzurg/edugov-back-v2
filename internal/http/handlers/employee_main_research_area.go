package handlers

import (
	"edugov-back-v2/internal/apperr"
	"edugov-back-v2/internal/dto"
	"edugov-back-v2/internal/http/middleware"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handlers) CreateEmployeeMainResearchArea(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEmployeeMainResearchAreaRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())
	req.LanguageCode = lang

	// DTO requires KeyTopics + each has LanguageCode required; also service ignores it,
	// so we enforce it here to keep requests consistent.
	if len(req.KeyTopics) == 0 {
		WriteError(w, r, apperr.Validation("validation failed", map[string]string{
			"keyTopics": "required",
		}))
		return
	}
	for i := range req.KeyTopics {
		if req.KeyTopics[i] == nil {
			WriteError(w, r, apperr.Validation("validation failed", map[string]string{
				"keyTopics[" + strconv.Itoa(i) + "]": "required",
			}))
			return
		}

		req.KeyTopics[i].LanguageCode = lang
		if req.KeyTopics[i].KeyTopicTitle == "" {
			WriteError(w, r, apperr.Validation("validation failed", map[string]string{
				"keyTopics[" + strconv.Itoa(i) + "].keyTopicTitle": "required",
			}))
			return
		}
	}

	resp, err := h.service.CreateEmployeeMainResearchArea(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handlers) GetEmployeeMainResearchAreas(w http.ResponseWriter, r *http.Request) {
	employeeID, err := h.readInt64QueryParam(r, "employeeID")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	lang := middleware.MustLang(r.Context())

	resp, err := h.service.GetEmployeeMainResearchAreaByEmployeeIDAndLanguageCode(r.Context(), employeeID, lang)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if resp == nil {
		resp = []*dto.EmployeeMainResearchAreaResponse{}
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) UpdateEmployeeMainResearchArea(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateEmployeeMainResearchAreaRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, r, err)
		return
	}

	for i := range req.KeyTopics {
		if req.KeyTopics[i] == nil {
			WriteError(w, r, apperr.Validation("validation failed", map[string]string{
				"keyTopics[" + strconv.Itoa(i) + "]": "required",
			}))
			return
		}

		if req.KeyTopics[i].ID < 0 {
			WriteError(w, r, apperr.Validation("validation failed", map[string]string{
				"keyTopics[" + strconv.Itoa(i) + "].id": "required",
			}))
			return
		}

		if req.KeyTopics[i].KeyTopicTitle == nil || strings.TrimSpace(*req.KeyTopics[i].KeyTopicTitle) == "" {
			WriteError(w, r, apperr.Validation("validation failed", map[string]string{
				"keyTopics[" + strconv.Itoa(i) + "].keyTopicTitle": "required",
			}))
			return
		}
	}

	resp, err := h.service.UpdateEmployeeMainResearchArea(r.Context(), &req)
	if err != nil {
		WriteError(w, r, err)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}

func (h *Handlers) DeleteEmployeeMainResearchArea(w http.ResponseWriter, r *http.Request) {
	id, err := h.readInt64URLParam(r, "id")
	if err != nil {
		WriteError(w, r, err)
		return
	}

	if err := h.service.DeleteEmployeeMainResearchArea(r.Context(), id); err != nil {
		WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
