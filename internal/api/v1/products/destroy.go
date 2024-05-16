package products

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/inconshreveable/log15"
)

func Destroy(w http.ResponseWriter, r *http.Request) {
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

	// Use access to the database to destroy the object
	if err := gr.Products().Destroy(r.Context(), id); err != nil {
		logger.Debug("unable to destroy product", log15.Ctx{"err": err, "id": id, "requestId": requestID})
		http.Error(w, "unable to destroy product id: "+requestID, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("success"))
}
