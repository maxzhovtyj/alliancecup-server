package inventory

import (
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

type Storage interface {
}

type storage struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func NewInventoryStorage(db *sqlx.DB, logger *logging.Logger) Storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}
