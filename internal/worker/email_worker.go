package worker

import (
	"context"
	"time"

	"github.com/dotuan0601/go-backend-template/internal/repository"
)

type EmailWorker struct {
	jobRepository repository.JobRepository

	workerID string
}

func NewEmailWorker(
	jobRepository repository.JobRepository,
	workerID string,
) *EmailWorker {

	return &EmailWorker{
		jobRepository: jobRepository,
		workerID:      workerID,
	}
}

func (w *EmailWorker) Start(
	ctx context.Context,
) {

	ticker := time.NewTicker(
		5 * time.Second,
	)

	go func() {

		for {

			select {

			case <-ctx.Done():

				return

			case <-ticker.C:

				w.process(ctx)

			}
		}

	}()

}

func (w *EmailWorker) process(
	ctx context.Context,
) {

	jobs, err := w.jobRepository.ClaimJobs(
		ctx,
		w.workerID,
		200,
	)

	if err != nil {
		return
	}

	for _, job := range jobs {

		// fake send email

		err := w.jobRepository.MarkSent(
			ctx,
			job.ID,
		)

		if err != nil {

			w.jobRepository.MarkFailed(
				ctx,
				job.ID,
				err,
			)
		}
	}
}
