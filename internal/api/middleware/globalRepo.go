package middleware

import (
	"context"
	"net/http"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"

	"github.com/gorilla/mux"
)

const ContextGlobalRepoKey contextKey = "mw:GlobalRepo"

func SetGlobalRepoOnContext(gr repos.GlobalRepo, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextGlobalRepoKey, gr))
}

func InjectGlobalRepo(gr repos.GlobalRepo) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextGlobalRepoKey, gr)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RetrieveGlobalRepo(ctx context.Context) (repos.GlobalRepo, bool) {
	ctxVal := ctx.Value(ContextGlobalRepoKey)
	val, exists := ctxVal.(repos.GlobalRepo)
	return val, exists
}
