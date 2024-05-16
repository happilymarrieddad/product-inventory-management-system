package products_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/v1/products"
	mock_repos "github.com/happilymarrieddad/product-inventory-management-system/internal/repos/mocks"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("HTTP: /v1/products", func() {
	var (
		ctrl         *gomock.Controller
		mockGr       *mock_repos.MockGlobalRepo
		mockProducts *mock_repos.MockProducts
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		mockGr = mock_repos.NewMockGlobalRepo(ctrl)
		mockProducts = mock_repos.NewMockProducts(ctrl)

		mockGr.EXPECT().Products().Return(mockProducts).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("/v1/products POST - create", func() {
		var body []byte
		BeforeEach(func() {
			var err error
			body, err = json.Marshal(types.NewProduct{
				Name: "some name", Sku: "some sku", Qty: 50,
			})
			Expect(err).To(BeNil())
		})

		It("should return an error when an invalid body is passed in", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("POST", "/v1/products", nil),
			)
			w := httptest.NewRecorder()

			products.Create(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to read body"))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should sanitize the err from the repo", func() {
			err := types.NewBadRequestError("BOGUS:Products.create")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("POST", "/v1/products", bytes.NewBuffer(body)),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Create(gomock.Any(), types.NewProduct{
				Name: "some name", Sku: "some sku", Qty: 50,
			}).Return(nil, err).Times(1)

			products.Create(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to create product"))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should successfully create a product", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("POST", "/v1/products", bytes.NewBuffer(body)),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Create(gomock.Any(), types.NewProduct{
				Name: "some name", Sku: "some sku", Qty: 50,
			}).Return(&types.Product{
				Name: "some name", Sku: "some sku", Qty: 50,
			}, nil).Times(1)

			products.Create(w, req)

			resp := w.Result()

			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
		})
	})

	Context("/v1/products/<id> DELETE - destroy", func() {
		It("should successfully destroy a product")
	})
})
