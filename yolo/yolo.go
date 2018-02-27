package yolo

import (
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZanLabs/go-yolo"
	logging "github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

// Task represents a yolo task
type Task struct {
	cate       types.Category
	ngpu       int
	dataCfg    string
	cfgFile    string
	weightFile string
	inputFile  string
	outputFile string
	thresh     float64
	hierThresh float64
}

// StartTask starts a yolo task
func StartTask(yt *Task, logger *logging.Logger, dbInstance db.Storage, jobID string) {
	// Select a gpu device randomly
	rand.Seed(time.Now().UnixNano())
	gpu := rand.Intn(yt.ngpu - 1)

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		logger.Error("retrieving-job", err)
		return
	}
	job.Status = types.JobProcessing
	job, err = dbInstance.UpdateJob(job.ID, job)
	if err != nil {
		logger.Error("processing-job", err)
		return
	}

	if yt.cate == types.IMAGE {
		// set gpu device
		yolo.SetGPU(gpu)
		// detect a image
		yolo.ImageDetector(
			yt.dataCfg,
			yt.cfgFile,
			yt.weightFile,
			yt.inputFile,
			yt.thresh,
			yt.hierThresh,
			yt.outputFile)
	} else if yt.cate == types.VIDEO {
		// set gpu device
		yolo.SetGPU(gpu)
		// detect a video
		yolo.VideoDetector(
			yt.dataCfg,
			yt.cfgFile,
			yt.weightFile,
			yt.inputFile,
			yt.thresh,
			yt.hierThresh,
			yt.outputFile)
	}
}

// NewTask creates a new task.
func NewTask(yc types.YoloConfig, cate types.Category, ngpu int, input string, output string) *Task {
	return &Task{
		cate:       cate,
		ngpu:       ngpu,
		dataCfg:    yc.DataCfg,
		cfgFile:    yc.CfgFile,
		weightFile: yc.WeightFile,
		inputFile:  input,
		// strip the suffix of the output file
		outputFile: filepath.Join(filepath.Dir(output), strings.Split(filepath.Base(output), ".")[0]),
		thresh:     yc.Thresh,
		hierThresh: yc.HierThresh,
	}
}
