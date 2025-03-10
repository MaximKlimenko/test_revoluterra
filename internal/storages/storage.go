package storages

type Storage interface {
	CreateJob(j *Job) error
	GetJobs() ([]Job, error)
	GetJobByID(id string) (Job, error)
	CancelJob(id string) error
	UpdateStatus(id string)
}
