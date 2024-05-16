package products_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProducts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Products Suite")
}
