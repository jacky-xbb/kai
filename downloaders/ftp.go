package downloaders

import (
	"net/url"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/secsy/goftp"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// FTPDownload downloads the file from FTP. Job Source should be
// in format: ftp://login:password@host/path
func FTPDownload(logger *logging.Logger, config types.SystemConfig, dbInstance db.Storage, jobID string) error {
	logger.Infof("start-ftp-download jobid: %s", jobID)
	defer logger.Infof("end-ftp-download jobid: %s", jobID)

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

	u, err := url.Parse(job.Source)
	if err != nil {
		return err
	}

	pw, isSet := u.User.Password()
	if !isSet {
		pw = ""
	}

	ftpConfig := goftp.Config{
		User:               u.User.Username(),
		Password:           pw,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(ftpConfig, u.Host+":21")
	if err != nil {
		logger.Error("dial-config-failed", err)
		return err
	}

	outputFile, err := os.Create(job.LocalSource)
	if err != nil {
		logger.Error("creating-local-source-failed", err)
		return err
	}

	err = client.Retrieve(u.Path, outputFile)
	if err != nil {
		logger.Error("retrieving-output-failed", err)
		return err
	}

	return nil
}
