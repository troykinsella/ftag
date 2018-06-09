package tagmap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTagmap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tagmap Suite")
}
