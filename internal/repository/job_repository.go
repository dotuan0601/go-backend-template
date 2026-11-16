package repository

import (
	"context"
	"database/sql"

	"github.com/dotuan0601/go-backend-template/internal/model"
)

type JobRepository interface {
	Create(
		ctx context.Context,
		job *model.Job,
	) error

	ClaimJobs(
		ctx context.Context,
		workerID string,
		limit int,
	) ([]*model.Job, error)

	MarkSent(
		ctx context.Context,
		id int64,
	) error

	MarkFailed(
		ctx context.Context,
		id int64,
		err error,
	) error
}

type jobRepository struct {
	db *sql.DB
}

func NewJobRepository(
	db *sql.DB,
) JobRepository {

	return &jobRepository{
		db: db,
	}
}

func (r *jobRepository) Create(
	ctx context.Context,
	job *model.Job,
) error {

	query := `
		INSERT INTO jobs
		(
			type,
			payload,
			status
		)
		VALUES
		(
			?,
			?,
			?
		)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		job.Type,
		job.Payload,
		job.Status,
	)

	if err != nil {

		return err
	}

	id, err := result.LastInsertId()

	if err != nil {

		return err
	}

	job.ID = id

	return nil
}

func (r *jobRepository) ClaimJobs(
	ctx context.Context,
	workerID string,
	limit int,
) ([]*model.Job, error) {

	tx, err := r.db.BeginTx(
		ctx,
		nil,
	)

	if err != nil {

		return nil, err
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(
		ctx,
		`
		SELECT
			id,
			type,
			payload,
			status,
			retry_count,
			created_at,
			updated_at
		FROM jobs
		WHERE status = 'PENDING'
		ORDER BY id
		LIMIT ?
		FOR UPDATE SKIP LOCKED
		`,
		limit,
	)

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	var jobs []*model.Job

	for rows.Next() {

		job := &model.Job{}

		err := rows.Scan(
			&job.ID,
			&job.Type,
			&job.Payload,
			&job.Status,
			&job.RetryCount,
			&job.CreatedAt,
			&job.UpdatedAt,
		)

		if err != nil {

			return nil, err
		}

		jobs = append(
			jobs,
			job,
		)
	}

	for _, job := range jobs {

		_, err := tx.ExecContext(
			ctx,
			`
			UPDATE jobs
			SET
				status = 'PROCESSING',
				locked_by = ?,
				locked_at = NOW()
			WHERE id = ?
			`,
			workerID,
			job.ID,
		)

		if err != nil {

			return nil, err
		}

		job.Status = model.JobStatusProcessing
	}

	if err := tx.Commit(); err != nil {

		return nil, err
	}

	return jobs, nil
}

func (r *jobRepository) MarkSent(
	ctx context.Context,
	id int64,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE jobs
		SET
			status = 'SENT',
			locked_by = NULL,
			locked_at = NULL
		WHERE id = ?
		`,
		id,
	)

	return err
}

func (r *jobRepository) MarkFailed(
	ctx context.Context,
	id int64,
	jobErr error,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE jobs
		SET
			status = 'FAILED',
			retry_count = retry_count + 1,
			locked_by = NULL,
			locked_at = NULL
		WHERE id = ?
		`,
		id,
	)

	return err
}
