package products_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
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

	Context("/v1/products/<id> GET - get", func() {
		It("should return an error when the repo is not on the context", func() {
			req := httptest.NewRequest("GET", "/v1/products", nil)
			w := httptest.NewRecorder()

			products.Get(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get internal resources"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should return a url parsing error", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, httptest.NewRequest("GET", "/v1/products/1", nil),
			)
			w := httptest.NewRecorder()

			products.Get(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get id from url parameters"))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should sanitize the err from the repo when an internal error", func() {
			err := types.NewInternalServerError("BOGUS:Products.get")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("GET", "/v1/products/1", nil),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, err).Times(1)

			products.Get(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get product"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should sanitize the err from the repo when no item is found", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("GET", "/v1/products/1", nil),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Get(gomock.Any(), int64(1)).Return(nil, false, nil).Times(1)

			products.Get(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get product"))
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should successfully get a product", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("GET", "/v1/products/1", nil),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Get(gomock.Any(), int64(1)).Return(&types.Product{
				ID: 1, Name: "some product",
			}, true, nil).Times(1)

			products.Get(w, req)

			resp := w.Result()

			bts, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bts)).To(ContainSubstring("some product"))
		})
	})
})
