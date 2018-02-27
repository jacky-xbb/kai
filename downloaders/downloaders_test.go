package downloaders_test

import (
	"reflect"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/downloaders"
	"github.com/ZanLabs/kai/types"
)

func initLogger() *logging.Logger {
	backend := logging.InitForTesting(logging.DEBUG)
	logging.SetBackend(backend)
	return logging.MustGetLogger("test")
}

var _ = Describe("Downloaders", func() {
	var (
		logger     *logging.Logger
		dbInstance db.Storage
		downloader downloaders.DownloadFunc
		exampleJob types.Job
		cfg        types.SystemConfig
	)

	BeforeEach(func() {
		logger = initLogger()

		cfg = types.SystemConfig{DBDriver: "memory"}
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDB()
	})

	Context("GetDownloadFunc", func() {
		It("should return S3Download if source has amazonaws", func() {
			jobSource := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.jpg"
			downloadFunc := downloaders.GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/ZanLabs/kai/downloaders.S3Download"))
		})

		It("should return FTPDownload if source starts with ftp://", func() {
			jobSource := "ftp://login:password@host/source_here.jpg"
			downloadFunc := downloaders.GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/ZanLabs/kai/downloaders.FTPDownload"))
		})

		It("should return HTTPDownload if source starts with http://", func() {
			jobSource := "http://source_here.jpg"
			downloadFunc := downloaders.GetDownloadFunc(jobSource)
			funcName := runtime.FuncForPC(reflect.ValueOf(downloadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/ZanLabs/kai/downloaders.HTTPDownload"))
		})
	})

	runDownloadersSuite := func() {
		It("should return an error if source couldn't be fetched", func() {
			dbInstance.StoreJob(exampleJob)
			err := downloader(logger, cfg, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(
				ContainSubstring("no such host"),
				ContainSubstring("No filename could be determined"),
				ContainSubstring("no such file or directory"),
				ContainSubstring("The AWS Access Key Id you provided does not exist in our records")))
		})
	}

	Context("HTTP Downloader", func() {
		BeforeEach(func() {
			downloader = downloaders.HTTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "http://source_here.jpg",
				Destination: "s3://user@pass:/bucket/",
				Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDownloadersSuite()
	})

	Context("FTP Downloader", func() {
		BeforeEach(func() {
			downloader = downloaders.FTPDownload
			exampleJob = types.Job{
				ID:          "123",
				Source:      "ftp://login:password@host/source_here.jpg",
				Destination: "s3://user@pass:/bucket/",
				Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
				Status:      types.JobCreated,
				Details:     "",
			}
		})

		runDownloadersSuite()
	})

	Context("S3 Downloader", func() {
		BeforeEach(func() {
			downloader = downloaders.S3Download
			exampleJob = types.Job{
				ID:               "123",
				Source:           "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.jpg",
				Destination:      "s3://user@pass:/bucket/",
				Media:            types.MediaType{Cate: types.IMAGE, Container: "jpg"},
				Status:           types.JobCreated,
				Details:          "",
				LocalDestination: "/tmp/output_here.jpg",
			}
		})

		runDownloadersSuite()
	})
})
