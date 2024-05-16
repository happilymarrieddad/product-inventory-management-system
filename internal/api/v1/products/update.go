package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	"github.com/inconshreveable/log15"
)

func Update(w http.ResponseWriter, r *http.Request) {
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

	// Get the updated product fields from the body of the request
	body := new(types.UpdateProduct)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		logger.Debug("unable to read body", log15.Ctx{"err": err, "requestId": requestID})
		http.Error(w, "unable to read body id: "+requestID, http.StatusBadRequest)
		return
	}

	// ensure the id is what was used in the URL
	// normally here we'd do an authorization check but this is not
	// an authenticated API
	body.ID = id

	// Use access to the database to update the requested object
	newProduct, err := gr.Products().Update(r.Context(), body)
	if err != nil {
		logger.Debug("unable to update product", log15.Ctx{
			"err": err, "id": id, "requestId": requestID, "req": body,
		})
		http.Error(w, "unable to update product id: "+requestID, http.StatusInternalServerError)
		return
	}

	// Marshal back the response
	bts, err := json.Marshal(newProduct)
	if err != nil {
		logger.Debug("unable to marshal response back", log15.Ctx{"err": err, "requestId": requestID})
		// json package is heavily tested and this will never happen but we should check for it
		http.Error(w, "unable to marshal product id: "+requestID, http.StatusInternalServerError)
		return
	}

	w.Write(bts)
}
