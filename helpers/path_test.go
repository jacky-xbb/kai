package helpers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/helpers"
	"github.com/ZanLabs/kai/types"
)

var _ = Describe("path", func() {
	var (
		dbInstance db.Storage
		cfg        types.SystemConfig
	)

	BeforeEach(func() {
		cfg = types.SystemConfig{DBDriver: "memory", SwapDir: "/tmp/"}
		dbInstance, _ = db.GetDatabase(cfg)
	})

	AfterEach(func() {
		dbInstance.ClearDB()
	})

	It("GetLocalSourcPath should return the correct local source path based on job", func() {
		exampleJob := types.Job{
			ID:          "123",
			Source:      "http://www.example.com/image.jpg",
			Destination: "s3://user@pass:/bucket/",
			Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
			Status:      types.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(exampleJob)

		res, err := helpers.GetLocalSourcePath(cfg, exampleJob.ID)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("/tmp/123/src/"))

	})

	It("GetLocalDestination should return the correct local destination path based on job", func() {
		exampleJob := types.Job{
			ID:          "123",
			Source:      "http://www.example.com/image.jpg",
			Destination: "s3://user@pass:/bucket/",
			Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
			Status:      types.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(exampleJob)

		Expect(helpers.GetLocalDestination(cfg, dbInstance, exampleJob.ID)).To(Equal("/tmp/123/dst/image.jpg"))
	})

	It("GetOutputFilename should build output filename based on job", func() {
		exampleJob := types.Job{
			ID:          "123",
			Source:      "http://www.example.com/image.jpg",
			Destination: "s3://user@pass:/bucket/",
			Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
			Status:      types.JobCreated,
			Details:     "",
		}
		dbInstance.StoreJob(exampleJob)

		res, err := helpers.GetOutputFilename(dbInstance, exampleJob.ID)
		Expect(err).To(BeNil())
		Expect(res).To(Equal("image.jpg"))
	})
})
