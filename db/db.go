package db

import (
	"github.com/ZanLabs/kai/types"
)

// Storage defines functions for accessing data
type Storage interface {
	// Job methods
	StoreJob(types.Job) (types.Job, error)
	RetrieveJob(string) (types.Job, error)
	UpdateJob(string, types.Job) (types.Job, error)
	GetJobs() ([]types.Job, error)
	DeleteJob(string) (types.Job, error)
	ClearDB() error
}

// GetDatabase selects the correspond driver depend on config
func GetDatabase(config types.SystemConfig) (Storage, error) {
	driver := config.DBDriver
	if driver == "mongo" || driver == "mongodb" {
		return getMongoDB(config)
	}
	return getMemoryDB()
}
