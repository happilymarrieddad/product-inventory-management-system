package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/inconshreveable/log15"
)

func Get(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	gr, exists := middleware.RetrieveGlobalRepo(r.Context())
	if !exists {
		logger.Debug("unable to get global repo from context", log15.Ctx{"requestId": requestID})
		http.Error(w, "unable to get internal resources id: "+requestID, http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		logger.Debug("unable to get id from url parameters", log15.Ctx{"err": err, "vars": mux.Vars(r), "requestId": requestID})
		http.Error(w, "unable to get id from url parameters id: "+requestID, http.StatusBadRequest)
		return
	}

	// Use access to the database to find the requested object
	product, exists, err := gr.Products().Get(r.Context(), id)
	if err != nil {
		logger.Debug("unable to get product", log15.Ctx{"err": err, "id": id, "requestId": requestID})
		http.Error(w, "unable to get product id: "+requestID, http.StatusInternalServerError)
		return
	}
	if !exists {
		logger.Debug("unable to get product", log15.Ctx{"err": err, "id": id, "requestId": requestID})
		http.Error(w, "unable to get product id: "+requestID, http.StatusNotFound)
		return
	}

	// Marshal back the response
	bts, err := json.Marshal(product)
	if err != nil {
		logger.Debug("unable to marshal response back", log15.Ctx{"err": err, "requestId": requestID})
		// json package is heavily tested and this will never happen but we should check for it
		http.Error(w, "unable to marshal product id: "+requestID, http.StatusInternalServerError)
		return
	}

	w.Write(bts)
}
