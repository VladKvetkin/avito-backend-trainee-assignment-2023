package models

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

type CreateSegmentRequest struct {
	Name string `json:"name"`
}

type SegmentResultResponse struct {
	Result bool `json:"result"`
}

type DeleteSegmentRequest struct {
	Name string `json:"name"`
}

type ChangeSegmentRequest struct {
	AddSegments    []AddSegment `json:"add_segments"`
	DeleteSegments []string     `json:"delete_segments"`
	UserID         int          `json:"user_id"`
}

type AddSegment struct {
	Name string `json:"name"`
	TTL  string `json:"ttl,omitempty"`
}

type GetSegmentsRequest struct {
	UserID int `json:"user_id"`
}

type GetSegmentsResponse struct {
	Segments []SegmentResponse `json:"segments"`
}

type SegmentResponse struct {
	Name string `json:"name"`
}

type GetSegmentsHistoryReportRequest struct {
	Period string `json:"period"`
}

type GetSegmentsHistoryReportResponse struct {
	URL string `json:"url"`
}
