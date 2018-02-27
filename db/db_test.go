package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/types"
)

var _ = Describe("Db", func() {
	var (
		dbInstance db.Storage
		runDBSuite func()

		job types.Job
	)

	BeforeEach(func() {
		job = types.Job{
			ID:          "123",
			Source:      "http://source1.here.jpg",
			Destination: "s3://user@pass:/bucket/destination1.jpg",
			Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
			Status:      types.JobCreated,
			Details:     "0%",
		}
	})

	runDBSuite = func() {
		Describe("StoreJob", func() {
			It("should be able to store a job", func() {
				res, err := dbInstance.StoreJob(job)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(job))
			})
		})

		Describe("RetrieveJob", func() {
			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
			})

			It("should be able to retrieve a job by its name", func() {
				res, err := dbInstance.RetrieveJob("123")
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(job))
			})

			Context("when the job does not exist", func() {
				It("should return an error", func() {
					_, err := dbInstance.RetrieveJob("invalid-job")
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("GetJobs", func() {
			var anotherJob types.Job

			BeforeEach(func() {
				anotherJob = types.Job{
					ID:          "321",
					Source:      "http://source2.here.jpg",
					Destination: "s3://user@pass:/bucket/destination2.jpg",
					Media:       types.MediaType{Cate: types.IMAGE, Container: "jpg"},
					Status:      types.JobCreated,
					Details:     "0%",
				}
			})

			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
				dbInstance.StoreJob(anotherJob)
			})

			It("should be able to list jobs", func() {
				jobs, err := dbInstance.GetJobs()
				Expect(err).NotTo(HaveOccurred())
				Expect(jobs).To(ConsistOf(job, anotherJob))
			})
		})

		Describe("UpdateJob", func() {
			JustBeforeEach(func() {
				dbInstance.StoreJob(job)
			})

			It("should be able to update job", func() {
				expectedStatus := types.JobDownloading
				job.Status = expectedStatus
				dbInstance.UpdateJob(job.ID, job)

				res, err := dbInstance.GetJobs()
				Expect(err).NotTo(HaveOccurred())
				Expect(res[0].Status).To(Equal(expectedStatus))
			})
		})
	}

	Describe("when the storage is in memory", func() {
		BeforeEach(func() {
			cfg := types.SystemConfig{DBDriver: "memory"}
			dbInstance, _ = db.GetDatabase(cfg)
		})

		AfterEach(func() {
			dbInstance.ClearDB()
		})
		runDBSuite()
	})

	Describe("when the storage is mongodb", func() {
		Describe("When it connects", func() {
			BeforeEach(func() {
				cfg := types.SystemConfig{DBDriver: "mongo", MongoHost: "localhost"}
				dbInstance, _ = db.GetDatabase(cfg)
			})

			AfterEach(func() {
				dbInstance.ClearDB()
			})

			runDBSuite()
		})

		Describe("when it fail to connect", func() {
			It("should not connect on mongo", func() {
				failedCfg := types.SystemConfig{DBDriver: "mongo", MongoHost: "invalid.ip.address"}
				_, err := db.GetDatabase(failedCfg)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
