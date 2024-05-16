package repos_test

import (
	"fmt"
	"time"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("REPOS: Products", func() {

	var (
		repo repos.Products
	)

	BeforeEach(func() {
		clearDatabase("products")

		repo = gr.Products()
		Expect(repo).NotTo(BeNil())
	})

	Context("Create(Tx)", func() {
		It("should fail with invalid products", func() {
			_, err := repo.Create(ctx, types.NewProduct{})
			Expect(err).NotTo(BeNil())

			_, err = repo.Create(ctx, types.NewProduct{Name: "test"})
			Expect(err).NotTo(BeNil())
			_, err = repo.Create(ctx, types.NewProduct{Sku: "test"})
			Expect(err).NotTo(BeNil())
			_, err = repo.Create(ctx, types.NewProduct{Qty: 50})
			Expect(err).NotTo(BeNil())

			_, err = repo.Create(ctx, types.NewProduct{Name: "test", Sku: "test"})
			Expect(err).NotTo(BeNil())
			_, err = repo.Create(ctx, types.NewProduct{Qty: 50, Sku: "test"})
			Expect(err).NotTo(BeNil())
		})

		It("should successfully create a product", func() {
			start := time.Now()
			newProduct, err := repo.Create(ctx, types.NewProduct{Name: "test", Sku: "test", Qty: 50})
			Expect(err).To(BeNil())
			Expect(newProduct).NotTo(BeNil())

			Expect(newProduct.ID).To(BeNumerically(">", 0))
			Expect(newProduct.Name).To(Equal("test"))
			Expect(newProduct.CreatedAt.After(start)).To(BeTrue())
			Expect(time.Now().After(newProduct.CreatedAt)).To(BeTrue())
		})
	})

	Context("product data creation", func() {
		var ids []int64
		BeforeEach(func() {
			ids = []int64{} // reset
			// Create 10 for testing
			for i := 0; i < 10; i++ {
				newProduct, err := repo.Create(ctx, types.NewProduct{
					Name: fmt.Sprintf("test-%d", i), Sku: fmt.Sprintf("sku-%d", i), Qty: 50,
				})
				Expect(err).To(BeNil())
				Expect(newProduct).NotTo(BeNil())
				Expect(newProduct.ID).To(BeNumerically(">", 0))

				ids = append(ids, newProduct.ID)
			}
		})

		Context("Find(Tx)", func() {
			It("should successfully return the full list", func() {
				products, count, err := repo.Find(ctx, nil)
				Expect(err).To(BeNil())
				Expect(count).To(BeNumerically("==", 10))
				Expect(products).To(HaveLen(10))
				Expect(products[0].ID).To(Equal(ids[0]))
				Expect(products[9].ID).To(Equal(ids[9]))
			})

			It("should successfully paginate", func() {
				products, count, err := repo.Find(ctx, &repos.ProductsFind{Limit: 2, Offset: 2})
				Expect(err).To(BeNil())
				// Should be all in the system
				Expect(count).To(BeNumerically("==", 10))
				Expect(products).To(HaveLen(2))
			})

			It("should allow looking for specific products", func() {
				products, count, err := repo.Find(ctx, &repos.ProductsFind{Limit: 1, IDs: []int64{ids[0]}})
				Expect(err).To(BeNil())
				Expect(count).To(BeNumerically("==", 1))
				Expect(products).To(HaveLen(1))
				Expect(products[0].ID).To(Equal(ids[0]))
			})
		})

		Context("Get(Tx)", func() {
			It("should not find an invalid product without an err", func() {
				product, exists, err := repo.Get(ctx, 9999999)
				Expect(err).To(BeNil())
				Expect(exists).To(BeFalse())
				Expect(product).To(BeNil())
			})

			It("should find a specific product", func() {
				product, exists, err := repo.Get(ctx, ids[0])
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
				Expect(product).NotTo(BeNil())
				Expect(product.ID).To(Equal(ids[0]))
			})
		})

		Context("Update(Tx)", func() {
			It("should return an error when attempting to update a product that doees not exist", func() {
				newName := "Some New Name"
				newProduct, err := repo.Update(ctx, &types.UpdateProduct{ID: 99999999, Name: &newName})
				Expect(err).NotTo(BeNil())
				Expect(err).To(Equal(types.NewNotFoundError("product not found by id")))
				Expect(newProduct).To(BeNil())
			})

			It("should successfully update", func() {
				product, exists, err := repo.Get(ctx, ids[0])
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
				Expect(product).NotTo(BeNil())
				Expect(product.ID).To(Equal(ids[0]))

				newName := "Some New Name"
				newProduct, err := repo.Update(ctx, &types.UpdateProduct{ID: ids[0], Name: &newName})
				Expect(err).To(BeNil())

				Expect(newProduct.Name).NotTo(Equal(product.Name))
				Expect(newProduct.Name).To(Equal(newName))

				Expect(newProduct.Sku).To(Equal(product.Sku))
				Expect(newProduct.Qty).To(Equal(product.Qty))
			})
		})

		Context("Destroy(Tx)", func() {
			It("should return an error when attempting to delete product that doesn't exist", func() {
				Expect(repo.Destroy(ctx, 999999999)).NotTo(Succeed())
			})

			It("should successfully delete a product that exists", func() {
				products, count, err := repo.Find(ctx, nil)
				Expect(err).To(BeNil())
				Expect(count).To(BeNumerically("==", 10))
				Expect(products).To(HaveLen(10))
				Expect(products[0].ID).To(Equal(ids[0]))
				Expect(products[9].ID).To(Equal(ids[9]))

				Expect(repo.Destroy(ctx, ids[0])).To(Succeed())

				products, count, err = repo.Find(ctx, nil)
				Expect(err).To(BeNil())
				Expect(count).To(BeNumerically("==", 9))
				Expect(products).To(HaveLen(9))
				Expect(products[0].ID).To(Equal(ids[1]))
				Expect(products[8].ID).To(Equal(ids[9]))
			})
		})
	})
})
