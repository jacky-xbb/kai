package downloaders

import (
	"os"

	"github.com/op/go-logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/helpers"
	"github.com/ZanLabs/kai/types"
)

// S3Download downloads the file from S3 bucket. Job Source should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Download(logger *logging.Logger, config types.SystemConfig, dbInstance db.Storage, jobID string) error {
	logger.Infof("start-s3-download jobid: %s", jobID)
	defer logger.Infof("end-s3-download jobid: %s", jobID)

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

	file, err := os.Create(job.LocalSource)
	if err != nil {
		return err
	}
	defer file.Close()

	err = helpers.SetAWSCredentials(job.Source)
	if err != nil {
		return err
	}

	bucket, err := helpers.GetAWSBucket(job.Source)
	if err != nil {
		return err
	}

	key, err := helpers.GetAWSKey(job.Source)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	objInput := s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}

	_, err = downloader.Download(file, &objInput)

	return err
}
