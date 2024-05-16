package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api/middleware"
	v1 "github.com/happilymarrieddad/product-inventory-management-system/internal/api/v1"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"
	"github.com/inconshreveable/log15"
)

var logger = log15.New("api")

func StartServer(port int, gr repos.GlobalRepo) {
	r := mux.NewRouter().StrictSlash(true)

	// Inject access to the database
	r.Use(middleware.InjectGlobalRepo(gr))

	// Add V1 routes
	v1.SetRoutes(r.PathPrefix("/v1").Subrouter())

	// Add handlers for error handling and methods|headers
	handler := handlers.LoggingHandler(os.Stdout, handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Origin", "Cache-Control", "X-App-Token"}),
		handlers.ExposedHeaders([]string{""}),
		handlers.MaxAge(1000),
		handlers.AllowCredentials(),
	)(r))
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)

	// Create the server with some reasonable timeouts
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	// Start server and inform the user what port the server is running on
	// note: log15 info complains with odd number of arguments...
	logger.Info("Server running", log15.Ctx{"port": port, "versions": []string{"v1"}})
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
