package uploaders

import (
	"os"

	"github.com/op/go-logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/helpers"
	"github.com/ZanLabs/kai/types"
)

// S3Upload sends the file to S3 bucket. Job Destination should be
// in format: http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/OBJECT
func S3Upload(logger *logging.Logger, dbInstance db.Storage, jobID string) error {
	logger.Infof("start-s3-upload jobid: %s", jobID)
	defer logger.Infof("end-s3-upload jobid: %s", jobID)

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}

	job.Status = types.JobUploading
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		logger.Error("downloading-job", err)
		return err
	}

	file, err := os.Open(job.LocalDestination)
	if err != nil {
		return err
	}

	err = helpers.SetAWSCredentials(job.Destination)
	if err != nil {
		return err
	}

	bucket, err := helpers.GetAWSBucket(job.Destination)
	if err != nil {
		return err
	}

	key, err := helpers.GetAWSKey(job.Destination)
	if err != nil {
		return err
	}

	job.Status = types.JobUploading
	dbInstance.UpdateJob(job.ID, job)

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
