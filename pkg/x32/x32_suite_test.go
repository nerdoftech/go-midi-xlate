package x32

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestX32(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "X32 Suite")
}
