package uploaders_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUploaders(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Uploaders Suite")
}
