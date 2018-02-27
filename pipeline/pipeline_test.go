package pipeline_test

import (
	"io"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	logging "github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/pipeline"
	"github.com/ZanLabs/kai/types"
)

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func createFakeFile() *os.File {
	f, err := ioutil.TempFile(os.TempDir(), "fake")
	if err != nil {
		panic(err)
	}
	return f
}

func initLogger() *logging.Logger {
	backend := logging.InitForTesting(logging.DEBUG)
	logging.SetBackend(backend)
	return logging.MustGetLogger("test")
}

var _ = Describe("Pipeline", func() {
	var (
		cfgSystem  types.SystemConfig
		logger     *logging.Logger
		dbInstance db.Storage
	)

	BeforeEach(func() {
		cfgSystem = types.SystemConfig{SwapDir: "/tmp/", DBDriver: "memory"}
		dbInstance, _ = db.GetDatabase(cfgSystem)
		logger = initLogger()
	})

	AfterEach(func() {
		dbInstance.ClearDB()
	})

	Context("SetupJob function", func() {
		It("Should set the local source and local destination on Job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://www.example.com/image.jpg",
				Destination: "s3://user@pass:/bucket/",
				Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
				Status:      types.JobCreated,
				Details:     "",
			}

			dbInstance.StoreJob(exampleJob)
			pipeline.SetupJob(logger, exampleJob.ID, dbInstance, cfgSystem)
			changedJob, _ := dbInstance.RetrieveJob("123")

			swapDir := cfgSystem.SwapDir

			sourceExpected := swapDir + "123/src/image.jpg"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/image.jpg"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	})

	Context("when calling Swap Cleaner", func() {
		It("should remove local source and local destination", func() {
			exampleJob := types.Job{
				ID:               "123",
				Source:           "http://source.here.jpg",
				Destination:      "s3://user@pass:/bucket/",
				Media:            types.MediaType{Cate: types.IMAGE, Container: "jpg"},
				Status:           types.JobCreated,
				Details:          "",
				LocalSource:      "/tmp/123/src/KailuaBeach.jpg",
				LocalDestination: "/tmp/123/dst/KailuaBeach.webm",
			}

			dbInstance.StoreJob(exampleJob)

			os.MkdirAll("/tmp/123/src/", 0777)
			os.MkdirAll("/tmp/123/dst/", 0777)

			fakeFile := createFakeFile()
			defer os.Remove(fakeFile.Name())

			cp(exampleJob.LocalSource, fakeFile.Name())
			cp(exampleJob.LocalDestination, fakeFile.Name())

			Expect(exampleJob.LocalSource).To(BeAnExistingFile())
			Expect(exampleJob.LocalDestination).To(BeAnExistingFile())

			pipeline.CleanSwap(dbInstance, exampleJob.ID)

			Expect(exampleJob.LocalSource).To(Not(BeAnExistingFile()))
			Expect(exampleJob.LocalDestination).To(Not(BeAnExistingFile()))
		})
	})
})
