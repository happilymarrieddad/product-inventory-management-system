package products

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	"github.com/inconshreveable/log15"
)

func Create(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	gr, exists := middleware.RetrieveGlobalRepo(r.Context())
	if !exists {
		logger.Debug("unable to get global repo from context", log15.Ctx{"requestId": requestID})
		http.Error(w, "unable to get internal resources id: "+requestID, http.StatusInternalServerError)
		return
	}

	// Get the new product from the body of the request
	body := new(types.NewProduct)
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		logger.Debug("unable to read body", log15.Ctx{"err": err, "requestId": requestID})
		http.Error(w, "unable to read body id: "+requestID, http.StatusBadRequest)
		return
	}

	// Use access to the database to create the new object
	np, err := gr.Products().Create(r.Context(), *body)
	if err != nil {
		logger.Debug("unable to read body", log15.Ctx{"err": err, "requestId": requestID})
		if types.IsBadRequestError(err) {
			http.Error(w, "unable to create product id: "+requestID, http.StatusBadRequest)
			return
		}
		http.Error(w, "unable to create product id: "+requestID, http.StatusInternalServerError)
		return
	}

	// Marshal back the response
	bts, err := json.Marshal(np)
	if err != nil {
		logger.Debug("unable to marshal response back", log15.Ctx{"err": err, "requestId": requestID})
		// json package is heavily tested and this will never happen but we should check for it
		http.Error(w, "unable to marshal product id: "+requestID, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(bts)
}
