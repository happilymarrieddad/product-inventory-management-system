package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"
	"github.com/happilymarrieddad/product-inventory-management-system/types"
	"github.com/inconshreveable/log15"
)

func Find(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	gr, exists := middleware.RetrieveGlobalRepo(r.Context())
	if !exists {
		logger.Debug("unable to get global repo from context", log15.Ctx{"requestId": requestID})
		http.Error(w, "unable to get internal resources id: "+requestID, http.StatusInternalServerError)
		return
	}

	opts := new(repos.ProductsFind)
	qry := r.URL.Query()

	limitRaw, exists := qry["limit"]
	if exists {
		limit, err := strconv.ParseInt(limitRaw[0], 10, 64)
		if err == nil {
			opts.Limit = int(limit)
		}
	}

	offsetRaw, exists := qry["offset"]
	if exists {
		offset, err := strconv.ParseInt(offsetRaw[0], 10, 64)
		if err == nil {
			opts.Offset = int(offset)
		}
	}

	idsRaw, exists := qry["id"]
	if exists {
		for _, idRaw := range idsRaw {
			id, err := strconv.ParseInt(idRaw, 10, 64)
			if err == nil {
				opts.IDs = append(opts.IDs, id)
			}
		}
	}

	nameRaw, exists := qry["name"]
	if exists {
		opts.Names = append(opts.Names, nameRaw...)
	}

	skuRaw, exists := qry["sku"]
	if exists {
		opts.Skus = append(opts.Skus, skuRaw...)
	}

	// Use access to the database to find the requested object(s)
	res, count, err := gr.Products().Find(r.Context(), opts)
	if err != nil {
		logger.Debug("unable to find products", log15.Ctx{"err": err, "requestId": requestID})
		if types.IsBadRequestError(err) {
			http.Error(w, "unable to find product id: "+requestID, http.StatusBadRequest)
			return
		}
		http.Error(w, "unable to find product id: "+requestID, http.StatusInternalServerError)
		return
	}

	// Marshal back the response
	bts, err := json.Marshal(struct {
		Data  interface{} `json:"data"`
		Count int64       `json:"count"`
	}{
		Data: res, Count: count,
	})
	if err != nil {
		logger.Debug("unable to marshal response back", log15.Ctx{"err": err, "requestId": requestID})
		// json package is heavily tested and this will never happen but we should check for it
		http.Error(w, "unable to marshal products id: "+requestID, http.StatusInternalServerError)
		return
	}

	w.Write(bts)
}
