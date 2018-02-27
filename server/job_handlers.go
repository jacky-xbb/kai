package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"

	"github.com/ZanLabs/kai/helpers"
	"github.com/ZanLabs/kai/pipeline"
	"github.com/ZanLabs/kai/types"
)

// CreateJob creates a job
func (ks *KaiServer) CreateJob(w http.ResponseWriter, r *http.Request) {
	ks.logger.Info("create-job")

	var jobInput types.JobInput
	if err := json.NewDecoder(r.Body).Decode(&jobInput); err != nil {
		ks.logger.Error("failed-unpacking-job", err)
		HTTPError(w, http.StatusBadRequest, "unpacking job", err)
		return
	}

	var job types.Job

	job.ID = uniuri.New()
	job.Source = jobInput.Source
	name, container, err := helpers.GetNameAndContainer(job.Source)
	if err != nil {
		ks.logger.Error("failed-extract-source-file", err)
		HTTPError(w, http.StatusBadRequest, "extracting source file", err)
		return
	}
	job.Media.Name = name
	job.Media.Container = container
	job.Media.Cate = jobInput.Cate
	job.Destination = jobInput.Destination
	job.Status = types.JobCreated

	_, err = ks.db.StoreJob(job)
	if err != nil {
		ks.logger.Error("failed-storing-job", err)
		HTTPError(w, http.StatusBadRequest, "storing job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		ks.logger.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", result)
	ks.logger.Infof("created jobid: %s", job.ID)
}

// StartJob triggers a recognition process
func (ks *KaiServer) StartJob(w http.ResponseWriter, r *http.Request) {
	ks.logger.Info("start-job")

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := ks.db.RetrieveJob(jobID)
	if err != nil {
		ks.logger.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	ks.logger.Debugf("starting-job id: %s", job.ID)
	// Todo: should return 202?
	w.WriteHeader(http.StatusOK)
	pipeline.Pipeline(ks.logger, ks.config, ks.db, ks.jobDownBuff, ks.jobTodoBuff, ks.jobDoneBuff, job)
}

// ListJobs lists all jobs
func (ks *KaiServer) ListJobs(w http.ResponseWriter, r *http.Request) {
	ks.logger.Info("list-jobs")

	jobs, err := ks.db.GetJobs()
	if err != nil {
		ks.logger.Error("failed-getting-jobs", err)
		HTTPError(w, http.StatusBadRequest, "getting jobs", err)
		return
	}

	result, err := json.Marshal(jobs)
	if err != nil {
		ks.logger.Error("failed-packaging-jobs", err)
		HTTPError(w, http.StatusBadRequest, "packing jobs data", err)
		return
	}

	fmt.Fprintf(w, "%s", string(result))
}

// GetJobDetails returns the details of a given job
func (ks *KaiServer) GetJobDetails(w http.ResponseWriter, r *http.Request) {
	ks.logger.Info("get-job-details")

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	job, err := ks.db.RetrieveJob(jobID)
	if err != nil {
		ks.logger.Error("failed-retrieving-job", err)
		HTTPError(w, http.StatusBadRequest, "retrieving job", err)
		return
	}

	result, err := json.Marshal(job)
	if err != nil {
		ks.logger.Error("failed-packaging-job-data", err)
		HTTPError(w, http.StatusBadRequest, "packing job data", err)
		return
	}

	fmt.Fprintf(w, "%s", result)
	ks.logger.Infof("got-job-details job: %+v", job)
}

// DeleteJob deletes a job
func (ks *KaiServer) DeleteJob(w http.ResponseWriter, r *http.Request) {
	ks.logger.Info("delete-job")

	vars := mux.Vars(r)
	jobID := vars["jobID"]
	_, err := ks.db.DeleteJob(jobID)
	if err != nil {
		ks.logger.Error("failed-deleting-job", err)
		HTTPError(w, http.StatusBadRequest, "deleting job", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
