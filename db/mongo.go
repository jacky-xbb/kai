package db

import (
	"sync"

	"github.com/ZanLabs/kai/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database struct that persists configurations
type mongoDB struct {
	db *mgo.Database
}

var databaseMongoInit sync.Once
var mongoInstance *mongoDB

// getMongoDB returns database singleton
func getMongoDB(config types.SystemConfig) (Storage, error) {
	var initErr error
	databaseMongoInit.Do(func() {
		mongoInstance = &mongoDB{}
	})

	mongoHost := config.MongoHost
	session, err := mgo.Dial(mongoHost)
	if err != nil {
		return nil, err
	}

	if mongoInstance.db == nil {
		session.SetMode(mgo.Monotonic, true)
		mongoInstance.db = session.DB("kai")
	}

	return mongoInstance, initErr
}

// ClearDatabase clears the database
func (r *mongoDB) ClearDatabase() error {
	return r.db.DropDatabase()
}

// StoreJob stores job information
func (r *mongoDB) StoreJob(job types.Job) (types.Job, error) {
	c := r.db.C("jobs")
	err := c.Insert(job)
	if err != nil {
		return types.Job{}, err
	}
	return job, nil
}

// RetrieveJob retrieves one job from the database
func (r *mongoDB) RetrieveJob(jobID string) (types.Job, error) {
	c := r.db.C("jobs")
	result := types.Job{}
	err := c.Find(bson.M{"id": jobID}).One(&result)
	return result, err
}

// UpdateJob updates a job
func (r *mongoDB) UpdateJob(jobID string, newJob types.Job) (types.Job, error) {
	c := r.db.C("jobs")
	err := c.Update(bson.M{"id": jobID}, newJob)
	if err != nil {
		return types.Job{}, err
	}
	return newJob, nil
}

//GetJobs retrieves all jobs of the database
func (r *mongoDB) GetJobs() ([]types.Job, error) {
	results := []types.Job{}
	c := r.db.C("jobs")
	err := c.Find(nil).All(&results)
	return results, err
}

// DeleteJob deletes a job from the database
func (r *mongoDB) DeleteJob(jobID string) (types.Job, error) {
	result, err := r.RetrieveJob(jobID)
	if err != nil {
		return types.Job{}, err
	}

	c := r.db.C("jobs")
	err = c.Remove(bson.M{"id": jobID})
	if err != nil {
		return types.Job{}, err
	}
	return result, nil
}

// ClearDB clears the database
func (r *mongoDB) ClearDB() error {
	return r.db.DropDatabase()
}
