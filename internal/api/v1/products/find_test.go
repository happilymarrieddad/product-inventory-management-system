package products_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/v1/products"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"
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

	Context("/v1/products GET - find", func() {
		It("should return an error when the repo is not on the context", func() {
			req := httptest.NewRequest("GET", "/v1/products", nil)
			w := httptest.NewRecorder()

			products.Find(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get internal resources"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should sanitize the err from the repo", func() {
			err := types.NewBadRequestError("BOGUS:Products.find")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("POST", "/v1/products", nil),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Find(gomock.Any(), gomock.AssignableToTypeOf(&repos.ProductsFind{})).
				Return(nil, int64(0), err).Times(1)

			products.Find(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to find product"))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should sanitize the err from the repo when an internal error", func() {
			err := types.NewInternalServerError("BOGUS:Products.find")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("POST", "/v1/products", nil),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Find(gomock.Any(), gomock.AssignableToTypeOf(&repos.ProductsFind{})).
				Return(nil, int64(0), err).Times(1)

			products.Find(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to find product"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should successfully find the products", func() {
			params := url.Values{}
			params.Add("limit", "25")
			params.Add("offset", "25")
			params.Add("id", "1234")
			params.Add("name", "1234")
			params.Add("sku", "1234")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("GET", "/v1/products", nil),
			)
			req.URL.RawQuery = "/v1/products?" + params.Encode()
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Find(gomock.Any(), &repos.ProductsFind{
				Limit: 25, Offset: 25, Names: []string{"1234"}, Skus: []string{"1234"},
			}).Return([]*types.Product{
				{Name: "some name", Sku: "some sku", Qty: 50},
				{Name: "some name 2", Sku: "some sku 2", Qty: 50},
			}, int64(2), nil).Times(1)

			products.Find(w, req)

			resp := w.Result()

			bts, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bts)).To(ContainSubstring("some name"))
			Expect(string(bts)).To(ContainSubstring("some name 2"))
		})
	})
})
