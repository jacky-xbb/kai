package downloaders

import (
	"strings"

	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// DownloadFunc is a function type for the multiple
// possible ways to download the source file
type DownloadFunc func(logger *logging.Logger, config types.SystemConfig, dbInstance db.Storage, jobID string) error

// GetDownloadFunc returns the download function
// based on the job source.
func GetDownloadFunc(jobSource string) DownloadFunc {
	if strings.Contains(jobSource, "aws") {
		return S3Download
	} else if strings.HasPrefix(jobSource, "ftp://") {
		return FTPDownload
	}

	return HTTPDownload
}
