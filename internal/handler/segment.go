package handler

import (
	"encoding/json"
	"net/http"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/models"
)

func (h *Handler) CreateSegment(res http.ResponseWriter, req *http.Request) {
	var createSegmentRequest models.CreateSegmentRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&createSegmentRequest); err != nil {
		h.errorJSONResponse(res, "cannot decode request to json", http.StatusBadRequest)
		return
	}

	if createSegmentRequest.Name == "" {
		h.errorJSONResponse(res, "cannot create segment with empty name", http.StatusBadRequest)
		return
	}

	if err := h.storage.CreateSegment(req.Context(), createSegmentRequest.Name); err != nil {
		h.errorJSONResponse(res, "error create segment", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(models.SegmentResultResponse{Result: true})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (h *Handler) DeleteSegment(res http.ResponseWriter, req *http.Request) {
	var deleteSegmentRequest models.DeleteSegmentRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&deleteSegmentRequest); err != nil {
		h.errorJSONResponse(res, "cannot decode request to json", http.StatusBadRequest)
		return
	}

	if deleteSegmentRequest.Name == "" {
		h.errorJSONResponse(res, "cannot delete segment with empty name", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteSegment(req.Context(), deleteSegmentRequest.Name); err != nil {
		h.errorJSONResponse(res, "error delete segment", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(models.SegmentResultResponse{Result: true})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
