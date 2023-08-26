package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/models"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/services/report"
)

func (h *Handler) ChangeSegments(res http.ResponseWriter, req *http.Request) {
	var changeSegmentsRequest models.ChangeSegmentRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&changeSegmentsRequest); err != nil {
		h.errorJSONResponse(res, "cannot decode request to json", http.StatusBadRequest)
		return
	}

	if len(changeSegmentsRequest.AddSegments) == 0 && len(changeSegmentsRequest.DeleteSegments) == 0 {
		h.errorJSONResponse(res, "empty change segments", http.StatusBadRequest)
		return
	}

	if err := h.storage.ChangeUserSegments(
		req.Context(),
		changeSegmentsRequest.UserID,
		changeSegmentsRequest.AddSegments,
		changeSegmentsRequest.DeleteSegments,
	); err != nil {
		h.errorJSONResponse(res, "error change user segments", http.StatusInternalServerError)
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

func (h *Handler) GetSegments(res http.ResponseWriter, req *http.Request) {
	var getSegmentsRequest models.GetSegmentsRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&getSegmentsRequest); err != nil {
		h.errorJSONResponse(res, "cannot decode request to json", http.StatusBadRequest)
		return
	}

	if getSegmentsRequest.UserID <= 0 {
		h.errorJSONResponse(res, "cannot get segment for empty user_id", http.StatusBadRequest)
		return
	}

	segments, err := h.storage.GetUserSegments(req.Context(), getSegmentsRequest.UserID)
	if err != nil {
		h.errorJSONResponse(res, "error get user segments from database", http.StatusInternalServerError)
		return
	}

	responseSegments := make([]models.SegmentResponse, 0, len(segments))
	for _, segment := range segments {
		respnonseSegment := models.SegmentResponse{
			Name: segment.Name,
		}

		responseSegments = append(responseSegments, respnonseSegment)
	}

	response, err := json.Marshal(models.GetSegmentsResponse{Segments: responseSegments})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (h *Handler) GetSegmentsHistoryReportURL(res http.ResponseWriter, req *http.Request) {
	var getSegmentsHistoryReportRequest models.GetSegmentsHistoryReportRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&getSegmentsHistoryReportRequest); err != nil {
		h.errorJSONResponse(res, "cannot decode request to json", http.StatusBadRequest)
		return
	}

	if getSegmentsHistoryReportRequest.Period == "" {
		h.errorJSONResponse(res, "empty segments history report period", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01", getSegmentsHistoryReportRequest.Period)
	if err != nil {
		h.errorJSONResponse(res, "invalid report period", http.StatusBadRequest)
		return
	}

	segmentsHistory, err := h.storage.GetSegmentsHistory(req.Context(), date)
	if err != nil {
		h.errorJSONResponse(res, "error get segments history from database", http.StatusInternalServerError)
		return
	}

	segmentsHistoryReport := report.NewCSVSegmentHistoryReport()

	reportURL, err := segmentsHistoryReport.GetReportURL(segmentsHistory, h.config.Address, getSegmentsHistoryReportRequest.Period)
	if err != nil {
		h.errorJSONResponse(res, "error create segments history report", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(models.GetSegmentsHistoryReportResponse{
		URL: reportURL,
	})

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (h *Handler) GetSegmentsHistoryReportFile(res http.ResponseWriter, req *http.Request) {
	reportFileHandler := http.StripPrefix("/api/segment/report", http.FileServer(http.Dir("./reports")))
	reportFileHandler.ServeHTTP(res, req)
}
