package postgres

import (
	"fmt"
	"time"

	"github.com/MaximKlimenko/scheduler/internal/config"
	"github.com/MaximKlimenko/scheduler/internal/storages"
)

type PostgresStorage struct {
	conn   *Connector
	Config *config.Config
}

func NewPostgresStorage(conn *Connector, config *config.Config) *PostgresStorage {
	return &PostgresStorage{
		conn:   conn,
		Config: config,
	}
}

func (s *PostgresStorage) CreateJob(j *storages.Job) error {
	return s.conn.DB.Create(j).Error
}

func (s *PostgresStorage) GetJobs() ([]storages.Job, error) {
	var jobs []storages.Job
	err := s.conn.DB.Find(&jobs).Error
	return jobs, err
}

func (s *PostgresStorage) GetJobByID(id string) (storages.Job, error) {
	var job storages.Job
	result := s.conn.DB.First(&job, "id = ?", id)
	if result.Error != nil {
		return storages.Job{}, result.Error
	}

	return job, nil
}

func (s *PostgresStorage) CancelJob(id string) error {
	var job storages.Job
	if err := s.conn.DB.First(&job, "id = ?", id).Error; err != nil {
		return err
	}

	if job.Status == storages.Cancelled || job.Status == storages.Executed {
		return fmt.Errorf("задача уже выполнена или отменена")
	}

	job.Status = storages.Cancelled
	s.conn.DB.Save(&job)
	return nil
}

func (s *PostgresStorage) UpdateStatus(id string) {
	var updatedJob storages.Job
	s.conn.DB.First(&updatedJob, "id = ?", id)
	if updatedJob.Status == storages.Scheduled {
		s.conn.DB.Model(&updatedJob).Update("status", storages.Executing)
		time.Sleep(2 * time.Second)
		executedAt := time.Now()
		s.conn.DB.Model(&updatedJob).Updates(storages.Job{Status: storages.Executed, ExecutedAt: &executedAt})
	}
}
