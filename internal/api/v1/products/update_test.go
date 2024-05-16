package products_test

import (
	"bytes"
	"encoding/json"
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

	Context("/v1/products PUT - update", func() {
		var body []byte
		BeforeEach(func() {
			var err error
			body, err = json.Marshal(types.NewProduct{
				Name: "some name", Sku: "some sku", Qty: 50,
			})
			Expect(err).To(BeNil())
		})

		It("should return an error when the repo is not on the context", func() {
			req := httptest.NewRequest("PUT", "/v1/products", nil)
			w := httptest.NewRecorder()

			products.Update(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to get internal resources"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should return an error when an invalid body is passed in", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("PUT", "/v1/products/1", nil),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			products.Update(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to read body"))
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})

		It("should sanitize the err from the repo", func() {
			err := types.NewBadRequestError("BOGUS:Products.update")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("PUT", "/v1/products/1", bytes.NewBuffer(body)),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(&types.UpdateProduct{})).Return(nil, err).Times(1)

			products.Update(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to update product"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should sanitize the err from the repo when an internal error", func() {
			err := types.NewInternalServerError("BOGUS:Products.update")

			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("PUT", "/v1/products/1", bytes.NewBuffer(body)),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(&types.UpdateProduct{})).Return(nil, err).Times(1)

			products.Update(w, req)

			resp := w.Result()

			resBts, _ := io.ReadAll(resp.Body)
			Expect(string(resBts)).To(ContainSubstring("unable to update product"))
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})

		It("should successfully update a product", func() {
			req := middleware.SetGlobalRepoOnContext(
				mockGr, mux.SetURLVars(httptest.NewRequest("PUT", "/v1/products/1", bytes.NewBuffer(body)),
					// Because of the helper function, we have to set it this way with gorilla mux
					map[string]string{"id": "1"},
				),
			)
			w := httptest.NewRecorder()

			mockProducts.EXPECT().Update(gomock.Any(), gomock.AssignableToTypeOf(&types.UpdateProduct{})).Return(&types.Product{
				Name: "some name", Sku: "some sku", Qty: 50,
			}, nil).Times(1)

			products.Update(w, req)

			resp := w.Result()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
