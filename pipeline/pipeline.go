package pipeline

import (
	"net/url"
	"os"
	"path"

	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/downloaders"
	"github.com/ZanLabs/kai/helpers"
	"github.com/ZanLabs/kai/types"
	"github.com/ZanLabs/kai/uploaders"
	"github.com/ZanLabs/kai/yolo"
)

// setupAndDownloadJob setup and download jobs into jobDownBuff
func setupAndDownloadJob(logger *logging.Logger, config types.SystemConfig,
	dbInstance db.Storage, job types.Job, jobDownBuff chan<- types.Job) {

	go func() {
		logger.Infof("start setup and download a job: %+v", job)
		newJob, err := SetupJob(logger, job.ID, dbInstance, config)
		job = *newJob
		if err != nil {
			logger.Error("setup-job failed", err)
			return
		}

		downloadFunc := downloaders.GetDownloadFunc(job.Source)
		if err := downloadFunc(logger, config, dbInstance, job.ID); err != nil {
			logger.Error("download failed", err)
			job.Status = types.JobError
			job.Details = err.Error()
			dbInstance.UpdateJob(job.ID, job)
			return
		}

		jobDownBuff <- job
	}()
}

func yoloJob(logger *logging.Logger, config types.ServerConfig, dbInstance db.Storage,
	jobDownBuff <-chan types.Job, jobTodoBuff chan types.Job, jobDoneBuff chan types.Job) {

	go func() {
		job, ok := <-jobDownBuff
		if !ok {
			logger.Info("job download buffer is closed")
			return
		}
		logger.Infof("start a yolo job: %+v", job)
		// limit the number of job in the jobTodoBuff
		jobTodoBuff <- job
		jobTodo, ok := <-jobTodoBuff
		if !ok {
			logger.Info("job todo buffer is closed")
			return
		}

		nGpu := config.System.NGpu
		t := yolo.NewTask(config.Yolo, jobTodo.Media.Cate, nGpu, jobTodo.LocalSource, jobTodo.LocalDestination)
		logger.Debugf("yolo task: %+v", *t)
		yolo.StartTask(t, logger, dbInstance, jobTodo.ID)
		jobDoneBuff <- job
	}()

}

func uploadJob(logger *logging.Logger, dbInstance db.Storage, jobDoneBuff <-chan types.Job) {
	go func() {
		jobDone, ok := <-jobDoneBuff
		if !ok {
			logger.Info("job done buffer is closed")
			return
		}
		logger.Infof("start a upload job: %+v", jobDone)

		uploadFunc := uploaders.GetUploadFunc(jobDone.Destination)
		if err := uploadFunc(logger, dbInstance, jobDone.ID); err != nil {
			logger.Error("upload failed", err)
			jobDone.Status = types.JobError
			jobDone.Details = err.Error()
			dbInstance.UpdateJob(jobDone.ID, jobDone)
			return
		}

		logger.Info("erasing temporary files")
		if err := CleanSwap(dbInstance, jobDone.ID); err != nil {
			logger.Error("erasing temporary files failed", err)
		}

		jobDone.Status = types.JobFinished
		dbInstance.UpdateJob(jobDone.ID, jobDone)

		logger.Infof("end a job: %+v", jobDone)
	}()
}

// Pipeline contains downloading, processing and uploading a job
func Pipeline(logger *logging.Logger, config types.ServerConfig, dbInstance db.Storage, jobDownBuff chan types.Job,
	jobTodoBuff chan types.Job, jobDoneBuff chan types.Job, job types.Job) {
	logger.Infof("pipeline-job %+v", job)

	// download a job
	setupAndDownloadJob(logger, config.System, dbInstance, job, jobDownBuff)

	// jobDownBuff -> jobTodoBuff -> jobDoneBuff
	yoloJob(logger, config, dbInstance, jobDownBuff, jobTodoBuff, jobDoneBuff)

	// upload a job
	uploadJob(logger, dbInstance, jobDoneBuff)
}

// CleanSwap removes LocalSource and LocalDestination
// files/directories.
func CleanSwap(dbInstance db.Storage, jobID string) error {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	err = os.RemoveAll(job.LocalSource)
	if err != nil {
		return err
	}

	err = os.RemoveAll(job.LocalDestination)
	return err
}

// SetupJob is responsible for set the initial state for a given
// job before starting. It sets local source and destination
// paths and the final destination as well.
func SetupJob(logger *logging.Logger, jobID string, dbInstance db.Storage, config types.SystemConfig) (*types.Job, error) {
	logger.Info("setup-job")
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return nil, err
	}

	localSource, err := helpers.GetLocalSourcePath(config, job.ID)
	if err != nil {
		return nil, err
	}
	job.LocalSource = localSource + path.Base(job.Source)
	logger.Debug("LocalSource: ", job.LocalSource)

	job.LocalDestination, err = helpers.GetLocalDestination(config, dbInstance, jobID)
	if err != nil {
		return nil, err
	}
	logger.Debug("LocalDestination: ", job.LocalDestination)

	u, err := url.Parse(job.Destination)
	if err != nil {
		return nil, err
	}
	outputFilename, err := helpers.GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, outputFilename)
	job.Destination = u.String()
	logger.Debug("Destination: ", job.Destination)
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
