package uploaders

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/op/go-logging"

	"github.com/secsy/goftp"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// FTPUpload uploades the file using FTP. Job Destination should be
// in format: ftp://login:password@host/path
func FTPUpload(logger *logging.Logger, dbInstance db.Storage, jobID string) error {
	logger.Infof("start-ftp-upload  jobid: %s", jobID)

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		logger.Error("retrieving-job", err)
		return err
	}

	job.Status = types.JobUploading
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		logger.Error("updating-job", err)
		return err
	}

	u, err := url.Parse(job.Destination)
	if err != nil {
		return err
	}

	pw, isSet := u.User.Password()
	if !isSet {
		pw = ""
	}

	config := goftp.Config{
		User:               u.User.Username(),
		Password:           pw,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
		Logger:             os.Stderr,
	}

	client, err := goftp.DialConfig(config, u.Host+":21")
	if err != nil {
		logger.Error("dial-config-failed", err)
		client.Close()
		return err
	}

	fileInfo, err := os.Stat(job.LocalDestination)
	if err != nil {
		logger.Error("get-destination-info", err)
		client.Close()
		return err
	}

	// remotePath := "." + path.Dir(u.Path)
	remotePath := path.Dir(u.Path)
	logger.Infof("check-remote-path path: %s", remotePath)

	_, err = client.Stat(remotePath)
	if err != nil {
		err = nil
		_, errMk := client.Mkdir(remotePath)
		if errMk != nil {
			logger.Error("no-create-path", errMk)
			client.Close()
			return errMk
		}
	}

	if fileInfo.IsDir() {
		base := path.Base(job.LocalDestination)
		client.Mkdir(remotePath + "/" + base)
		files, err := ioutil.ReadDir(job.LocalDestination)
		if err != nil {
			logger.Error("listing-files", err)
			client.Close()
			return err
		}
		for _, file := range files {
			localFile, err := os.Open(job.LocalDestination + "/" + file.Name())
			defer localFile.Close()
			if err != nil {
				logger.Error("opening-local-destination-failed", err)
				client.Close()
				return err
			}
			// client.Store("."+u.Path+"/"+file.Name(), localFile)
			client.Store(u.Path+"/"+file.Name(), localFile)
		}

	} else {
		localFile, err := os.Open(job.LocalDestination)
		defer localFile.Close()
		if err != nil {
			logger.Error("opening-local-destination-failed", err)
			client.Close()
			return err
		}

		// err = client.Store("."+u.Path, localFile)
		err = client.Store(u.Path, localFile)
		if err != nil {
			logger.Error("storing-file-failed", err)
			client.Close()
			return err
		}
	}
	client.Close()
	return err
}
