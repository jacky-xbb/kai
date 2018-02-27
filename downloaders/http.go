package downloaders

import (
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(logger *logging.Logger, config types.SystemConfig, dbInstance db.Storage, jobID string) error {
	logger.Infof("start-http-download jobid: %s", jobID)
	defer logger.Infof("end-http-download jobid: %s", jobID)

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	job.Status = types.JobDownloading
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		logger.Error("downloading-job", err)
		return err
	}
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(job.LocalSource, job.Source)

	// start download
	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			logger.Debug("transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(), resp.Size, 100*resp.Progress())

		case <-resp.Done:
			break Loop
		}
	}

	return resp.Err()
}
