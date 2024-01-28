package dbkeeper

import (
	"context"
	"database/sql"
)

type DBKeeper struct {
	db *sql.DB
}

func NewDBKeeper(db *sql.DB) *DBKeeper {
	return &DBKeeper{
		db: db,
	}
}

func (db *DBKeeper) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}
