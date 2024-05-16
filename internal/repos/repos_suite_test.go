package repos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/config"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/db"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var gr repos.GlobalRepo
var ctx context.Context

func clearDatabase(tables ...string) {
	for _, table := range tables {
		_, err := gr.DB().Exec(fmt.Sprintf("TRUNCATE %s CASCADE", table))
		Expect(err).To(Succeed())
	}
}

var _ = BeforeSuite(func() {
	defer GinkgoRecover()

	ctx = context.Background()
	cfg := config.NewConfig()
	cfg.Debug = true // force debug in testing

	var err error
	gr, err = db.NewDB(cfg.DBConfig)
	Expect(err).To(BeNil())
})

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
}
