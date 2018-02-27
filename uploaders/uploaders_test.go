package uploaders_test

import (
	"reflect"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ZanLabs/kai/uploaders"
)

var _ = Describe("Uploaders", func() {
	Context("GetUploadFunc", func() {
		It("should return S3Upload if source has amazonaws", func() {
			jobDestination := "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/source_here.jpg"
			uploadFunc := uploaders.GetUploadFunc(jobDestination)
			funcName := runtime.FuncForPC(reflect.ValueOf(uploadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/ZanLabs/kai/uploaders.S3Upload"))
		})

		It("should return FTPUpload if source starts with ftp://", func() {
			jobDestination := "ftp://login:password@host/source_here.jpg"
			uploadFunc := uploaders.GetUploadFunc(jobDestination)
			funcName := runtime.FuncForPC(reflect.ValueOf(uploadFunc).Pointer()).Name()
			Expect(funcName).To(Equal("github.com/ZanLabs/kai/uploaders.FTPUpload"))
		})
	})
})
