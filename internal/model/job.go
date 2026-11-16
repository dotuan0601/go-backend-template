package model

import "time"

type Job struct {
	ID int64

	Type string

	Payload string

	Status string

	RetryCount int

	LockedBy *string

	LockedAt *time.Time

	FailedAt *time.Time

	CreatedAt time.Time

	UpdatedAt time.Time
}
