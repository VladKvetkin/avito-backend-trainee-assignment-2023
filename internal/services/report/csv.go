package report

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/entities"
)

type SegmentHistoryReport interface {
	GetReportURL([]entities.SegmentHistory, string, string) (string, error)
}

type CSVSegmentHistoryReport struct {
}

func NewCSVSegmentHistoryReport() SegmentHistoryReport {
	return &CSVSegmentHistoryReport{}
}

func (report *CSVSegmentHistoryReport) GetReportURL(segmentsHistory []entities.SegmentHistory, address string, period string) (string, error) {
	csvFileName := fmt.Sprintf("%s-segments-history-report.csv", period)
	fullFileName := "reports/" + csvFileName

	if _, err := os.Stat(fullFileName); err == nil {
		if err := os.Remove(fullFileName); err != nil {
			return "", err
		}
	}

	file, err := os.Create(fullFileName)
	if err != nil {
		return "", err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'

	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"UserID", "Segment", "Operation", "Created"}); err != nil {
		return "", err
	}

	for _, segmentHistory := range segmentsHistory {
		if err := csvWriter.Write([]string{
			fmt.Sprintf("%d", segmentHistory.UserID),
			segmentHistory.SegmentName,
			segmentHistory.Operation,
			segmentHistory.CreatedAt.Format(time.DateTime),
		}); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/%s", "localhost:8080/api/segment/report", csvFileName), nil
}
