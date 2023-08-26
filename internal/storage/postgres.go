package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/entities"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/models"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	CreateSegment(context.Context, string) error
	DeleteSegment(context.Context, string) error
	GetUserSegments(context.Context, int) ([]entities.Segment, error)
	ChangeUserSegments(context.Context, int, []models.AddSegment, []string) error
	GetSegmentsHistory(context.Context, time.Time) ([]entities.SegmentHistory, error)
	DeleteExpiredUserSegments(context.Context) error
}

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(db *sqlx.DB) Storage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) CreateSegment(ctx context.Context, name string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO segment(name) VALUES ($1) ON CONFLICT DO NOTHING`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgresStorage) DeleteSegment(ctx context.Context, name string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM segment WHERE name = $1`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgresStorage) GetUserSegments(ctx context.Context, userID int) ([]entities.Segment, error) {
	var segments []entities.Segment

	if err := s.db.SelectContext(
		ctx,
		&segments,
		"SELECT segment_name AS name FROM user_segment WHERE user_id = $1 AND ttl > $2::timestamp ORDER BY name ASC;",
		userID,
		time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return nil, err
	}

	return segments, nil
}

func (s *PostgresStorage) ChangeUserSegments(ctx context.Context, userID int, addSegments []models.AddSegment, deleteSegments []string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, segment := range addSegments {
		s.CreateSegment(ctx, segment.Name)

		var result sql.Result

		if segment.TTL == "" {
			result, err = tx.ExecContext(
				ctx,
				`INSERT INTO user_segment (user_id, segment_name)
				 VALUES ($1, $2)
				 ON CONFLICT DO NOTHING`,
				userID,
				segment.Name,
			)

			if err != nil {
				return err
			}
		} else {
			ttl, err := time.Parse(time.DateTime, segment.TTL)
			if err != nil {
				return err
			}

			result, err = tx.ExecContext(
				ctx,
				`INSERT INTO user_segment (user_id, segment_name, ttl)
				 VALUES ($1, $2, $3::timestamp)
				 ON CONFLICT DO NOTHING`,
				userID,
				segment.Name,
				ttl.UTC().Format(time.RFC3339),
			)

			if err != nil {
				return err
			}
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected > 0 {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO user_segment_history (user_id, segment_name, operation)
				VALUES ($1, $2, $3)`,
				userID,
				segment.Name,
				entities.SegmentHistoryCreateOperation,
			); err != nil {
				return err
			}
		}
	}

	for _, segmentName := range deleteSegments {
		s.CreateSegment(ctx, segmentName)

		result, err := tx.ExecContext(
			ctx,
			`DELETE FROM user_segment WHERE user_id = $1 AND segment_name = $2`,
			userID,
			segmentName,
		)

		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected > 0 {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO user_segment_history (user_id, segment_name, operation)
				VALUES ($1, $2, $3)`,
				userID,
				segmentName,
				entities.SegmentHistoryDeleteOperation,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *PostgresStorage) GetSegmentsHistory(ctx context.Context, period time.Time) ([]entities.SegmentHistory, error) {
	var segmentsHistory []entities.SegmentHistory

	if err := s.db.SelectContext(
		ctx,
		&segmentsHistory,
		"SELECT * FROM user_segment_history WHERE created_at <= $1::timestamp ORDER BY created_at DESC;",
		period.UTC().Format(time.RFC3339),
	); err != nil {
		return nil, err
	}

	return segmentsHistory, nil
}

func (s *PostgresStorage) DeleteExpiredUserSegments(ctx context.Context) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_segment WHERE ttl < $1::timestamp`, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return err
	}

	return tx.Commit()
}
