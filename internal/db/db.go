package db

import (
	"fmt"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/repos"

	"github.com/happilymarrieddad/product-inventory-management-system/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"xorm.io/xorm"
)

// NewDB creates a new database connection
func NewDB(cfg config.DBConfig) (repos.GlobalRepo, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?connect_timeout=180&sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := xorm.NewEngine("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.ShowSQL(true)

	gr, err := repos.NewGlobalRepo(db)
	if err != nil {
		return nil, err
	}

	return gr, nil
}
