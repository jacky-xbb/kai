package db

import (
	"errors"
	"sync"

	"github.com/ZanLabs/kai/types"
)

// memoryDB struct that persists configurations
type memoryDB struct {
	mtx sync.RWMutex

	jobs map[string]types.Job
}

var dbInit sync.Once
var memoryInstance *memoryDB

// GetMemory returns database singleton
func getMemoryDB() (Storage, error) {
	dbInit.Do(func() {
		memoryInstance = &memoryDB{}
		memoryInstance.jobs = map[string]types.Job{}
	})

	return memoryInstance, nil
}

// CleanDatabase cleans the database
func (r *memoryDB) ClearDB() error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	memoryInstance.jobs = map[string]types.Job{}
	return nil
}

// RetrieveJob retrieves one job from the database
func (r *memoryDB) RetrieveJob(jobID string) (types.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.jobs[jobID]; ok {
		return val, nil
	}
	return types.Job{}, errors.New("job not found")
}

// StoreJob stores job information
func (r *memoryDB) StoreJob(job types.Job) (types.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[job.ID] = job
	return job, nil
}

// UpdateJob updates a job
func (r *memoryDB) UpdateJob(jobID string, newJob types.Job) (types.Job, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.jobs[jobID] = newJob
	return newJob, nil
}

// GetJobs retrieves all jobs of the database
func (r *memoryDB) GetJobs() ([]types.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	res := make([]types.Job, 0, len(r.jobs))
	for _, value := range r.jobs {
		res = append(res, value)
	}
	return res, nil
}

// DeleteJob delete the job with the jobID
func (r *memoryDB) DeleteJob(jobID string) (types.Job, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if val, ok := r.jobs[jobID]; ok {
		delete(r.jobs, jobID)
		return val, nil
	}
	return types.Job{}, errors.New("job not found")
}
