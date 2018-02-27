package server_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	logging "github.com/op/go-logging"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

func initLogger() *logging.Logger {
	backend := logging.InitForTesting(logging.DEBUG)
	logging.SetBackend(backend)
	return logging.MustGetLogger("test")
}
