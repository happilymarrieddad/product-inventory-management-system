package main

import (
	"github.com/happilymarrieddad/product-inventory-management-system/internal/api"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/config"
	"github.com/happilymarrieddad/product-inventory-management-system/internal/db"
)

func main() {
	// First, we grab the config from the env
	cfg := config.NewConfig()

	// Next, we grab access to the database using the config
	gr, err := db.NewDB(cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	// Now, we start the server
	api.StartServer(cfg.Port, gr)
}
