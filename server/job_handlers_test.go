package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/server"
	"github.com/ZanLabs/kai/types"
)

var _ = Describe("JobHandlers", func() {
	var (
		dbInstance       db.Storage
		ks               *server.KaiServer
		input            types.JobInput
		jobRecorder      *httptest.ResponseRecorder
		respJobInputBody map[string]interface{}
	)

	BeforeEach(func() {
		cfgSystem := types.SystemConfig{SwapDir: "/tmp/", DBDriver: "memory"}
		dbInstance, _ = db.GetDatabase(cfgSystem)

		logger := initLogger()

		cfgServer := types.ServerConfig{
			ServerName: "kai",
			System:     cfgSystem}
		ks = server.New(logger, cfgServer, "tcp", ":8000", dbInstance)

		input = types.JobInput{
			Source:      "http://s3.example.com/images/image1.jpg",
			Destination: "s3://example-bucket/future/image1.jpg",
			Cate:        types.IMAGE,
		}

		jobRecorder = httptest.NewRecorder()
		payloadData, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(payloadData))

		ks.Handler().ServeHTTP(jobRecorder, req)
		json.Unmarshal(jobRecorder.Body.Bytes(), &respJobInputBody)
	})

	AfterEach(func() {
		dbInstance.ClearDB()
	})

	Context("Create job", func() {
		It("should create a job in the db instance", func() {
			Expect(jobRecorder.Code).To(BeIdenticalTo(http.StatusCreated))
			jobID, ok := respJobInputBody["id"].(string)
			Expect(ok).To(BeIdenticalTo(true))
			job, err := dbInstance.RetrieveJob(jobID)
			Expect(err).NotTo(HaveOccurred())
			Expect(job.Destination).To(BeIdenticalTo(input.Destination))
		})

		It("should get a given job details", func() {
			var jobBody map[string]interface{}

			recorder := httptest.NewRecorder()
			reqJobDetail, _ := http.NewRequest(http.MethodGet, "/jobs/"+respJobInputBody["id"].(string), nil)
			ks.Handler().ServeHTTP(recorder, reqJobDetail)
			json.Unmarshal(recorder.Body.Bytes(), &jobBody)
			Expect(recorder.Code).To(BeIdenticalTo(http.StatusOK))
			Expect(jobBody["id"]).To(BeIdenticalTo(respJobInputBody["id"]))
			Expect(jobBody["status"]).To(BeIdenticalTo(respJobInputBody["status"]))
		})

		It("should list all jobs", func() {
			secondInput := types.JobInput{
				Source:      "http://s3.example.com/images/image2.jpg",
				Destination: "s3://example-bucket/future/image2.jpg",
				Cate:        types.IMAGE,
			}

			recorder := httptest.NewRecorder()
			data, _ := json.Marshal(secondInput)
			req, _ := http.NewRequest(http.MethodPost, "/jobs", bytes.NewReader(data))
			ks.Handler().ServeHTTP(recorder, req)

			listRecorder := httptest.NewRecorder()
			var jobListBody []map[string]interface{}

			listJobsRequest, _ := http.NewRequest(http.MethodGet, "/jobs", nil)
			ks.Handler().ServeHTTP(listRecorder, listJobsRequest)
			json.Unmarshal(listRecorder.Body.Bytes(), &jobListBody)
			Expect(listRecorder.Code).To(BeIdenticalTo(http.StatusOK))
			Expect(len(jobListBody)).To(Equal(2))
		})
	})
})
