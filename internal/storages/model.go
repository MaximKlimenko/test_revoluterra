package storages

import "time"

type JobStatus string

const (
	Scheduled JobStatus = "scheduled"
	Executing JobStatus = "executing"
	Executed  JobStatus = "executed"
	Cancelled JobStatus = "cancelled"
)

type Job struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Description string     `json:"description"`
	ExecuteAt   time.Time  `json:"executeAt"`
	Status      JobStatus  `json:"status"`
	ExecutedAt  *time.Time `json:"executedAt,omitempty"`
}
