package entities

import "time"

const (
	SegmentHistoryCreateOperation = "CREATE"
	SegmentHistoryDeleteOperation = "DELETE"
)

type SegmentHistory struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	SegmentName string    `db:"segment_name"`
	Operation   string    `db:"operation"`
	CreatedAt   time.Time `db:"created_at"`
}
