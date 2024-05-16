package products

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inconshreveable/log15"
)

var logger = log15.New("/v1/products")

func SetRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("", Find).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id:[0-9]+}", Get).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id:[0-9]+}", Update).Methods(http.MethodPut)
	subrouter.HandleFunc("/{id:[0-9]+}", Destroy).Methods(http.MethodDelete)
	subrouter.HandleFunc("", Create).Methods(http.MethodPost)
}
