package server_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/server"
	"github.com/ZanLabs/kai/types"
)

var _ = Describe("Kai Server", func() {
	var (
		logger     *logging.Logger
		ks         *server.KaiServer
		dbInstance db.Storage
		cfgServer  types.ServerConfig
	)

	BeforeEach(func() {
		logger = initLogger()

		cfgSystem := types.SystemConfig{SwapDir: "/tmp/", DBDriver: "memory"}
		dbInstance, _ = db.GetDatabase(cfgSystem)

		cfgServer = types.ServerConfig{
			ServerName: "kai",
			System:     cfgSystem}
	})

	AfterEach(func() {
		dbInstance.ClearDB()
	})

	Context("when passed a socket", func() {
		var (
			socketPath string
			tmpDir     string
		)

		JustBeforeEach(func() {
			var err error
			tmpDir, err = ioutil.TempDir(os.TempDir(), "kai-server-test")
			socketPath = path.Join(tmpDir, "kai.sock")
			ks = server.New(logger, cfgServer, "unix", socketPath, dbInstance)
			Expect(err).NotTo(HaveOccurred())

			err = ks.Start(false)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tmpDir)
		})

		Context("Start", func() {
			It("listens on the socket provided", func() {
				info, err := os.Stat(socketPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info).NotTo(BeNil())
			})
		})

		Context("Stop", func() {
			JustBeforeEach(func() {
				info, err := os.Stat(socketPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(info).NotTo(BeNil())
			})

			It("removes the existing socket", func() {
				Expect(ks.Stop()).To(Succeed())

				info, err := os.Stat(socketPath)
				Expect(err).To(HaveOccurred())
				Expect(info).To(BeNil())
			})

			Context("when fails to stop the server because it's already stopped", func() {
				JustBeforeEach(func() {
					Expect(ks.Stop()).NotTo(HaveOccurred())
				})

				It("returns an error", func() {
					Expect(ks.Stop()).To(HaveOccurred())
				})
			})
		})
	})
})
