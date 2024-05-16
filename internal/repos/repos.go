package repos

import (
	"log"
	"strings"

	"github.com/happilymarrieddad/product-inventory-management-system/types"
	"github.com/jackc/pgx/v5/pgconn"
	"xorm.io/xorm"
)

func handleRollback(sesh *xorm.Session, err error) error {
	if rollBackErr := sesh.Rollback(); rollBackErr != nil {
		log.Printf("unable to rollback with err: %s", rollBackErr.Error())
	}

	return err
}

func wrapInSession(db *xorm.Engine, fn func(*xorm.Session) (any, error)) (any, error) {
	session := db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return nil, types.NewInternalServerError("unable to start transaction with err: " + err.Error())
	}

	res, err := fn(session)
	if err != nil {
		return nil, handleRollback(session, err)
	}

	if err := session.Commit(); err != nil {
		return nil, handleRollback(session, err)
	}

	return res, nil
}

func normalizeErr(entityName string, err error) error {
	if err == nil {
		return nil
	}

	if err == xorm.ErrNotExist {
		return types.NewNotFoundError("unable to find attribute by id")
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if strings.Contains(pgErr.Message, "duplicate key value violates unique constraint") {
			log.Println("database err with constraint ", pgErr.ConstraintName)
			return types.NewBadRequestError("duplicate object already exists")
		}
	}

	return err
}
