package v1

import (
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/v1/products"
)

func SetRoutes(subrouter *mux.Router) {
	products.SetRoutes(subrouter.PathPrefix("/products").Subrouter())
}
