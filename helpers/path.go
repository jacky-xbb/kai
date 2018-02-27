package helpers

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// GetNameAndContainer returns the name and the container of a media file.
func GetNameAndContainer(source string) (string, string, error) {
	s1 := strings.Split(source, "/")
	s2 := strings.Split(s1[len(s1)-1], ".")
	if s2[0] == "" || s2[1] == "" {
		return "", "", errors.New("can't get name and container")
	}
	return s2[0], s2[1], nil
}

// GetLocalSourcePath builds the path and filename for
// the local source file
func GetLocalSourcePath(config types.SystemConfig, jobID string) (string, error) {
	baseDir, err := getBaseDir(config, jobID)
	sourceDir := baseDir + "/src/"
	if err != nil {
		return "", err
	}

	os.MkdirAll(sourceDir, 0700)

	return sourceDir, nil
}

// GetLocalDestination builds the path and filename
// of the local destination file
func GetLocalDestination(config types.SystemConfig, dbInstance db.Storage, jobID string) (string, error) {
	baseDir, err := getBaseDir(config, jobID)
	if err != nil {
		return "", err
	}

	destinationDir := baseDir + "/dst/"
	if err != nil {
		return "", err
	}

	os.MkdirAll(destinationDir, 0700)
	outputFilename, err := GetOutputFilename(dbInstance, jobID)
	if err != nil {
		return "", err
	}

	return destinationDir + outputFilename, nil
}

// GetOutputFilename build the destination path with
// the output filename
func GetOutputFilename(dbInstance db.Storage, jobID string) (string, error) {
	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return "", err
	}

	return strings.Split(path.Base(job.Source), ".")[0] + "." + job.Media.Container, nil
}

func getBaseDir(config types.SystemConfig, jobID string) (string, error) {
	swapDir := config.SwapDir

	return swapDir + jobID, nil
}
