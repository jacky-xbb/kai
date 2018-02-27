package uploaders

import (
	"strings"

	"github.com/ZanLabs/kai/db"
	logging "github.com/op/go-logging"
)

// UploadFunc is a function type for the multiple
// possible ways to upload the source file
type UploadFunc func(logger *logging.Logger, dbInstance db.Storage, jobID string) error

// GetUploadFunc returns the upload function
// based on the job source.
func GetUploadFunc(jobDestination string) UploadFunc {
	if strings.HasPrefix(jobDestination, "ftp://") {
		return FTPUpload
	}

	return S3Upload
}
