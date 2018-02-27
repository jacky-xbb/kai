package downloaders_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDownloaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Downloaders Suite")
}
