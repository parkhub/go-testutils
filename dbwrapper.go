package testutils

import (
	"github.com/go-pg/pg/v9"
)

type DBWrapper struct {
	*pg.DB
}

// RunInTransaction runs a function in a transaction. If function
// returns an error transaction is rollbacked, otherwise transaction
// is committed.
func (db *DBWrapper) RunInTransaction(fn func(Tx) error) error {
	var fn2 func(*pg.Tx) error
	fn2 = func(tx *pg.Tx) error {
		return fn(tx)
	}
	return db.DB.RunInTransaction(fn2)
}

// Query executes a query that returns rows, typically a SELECT.
// The params are for any placeholders in the query.
func (db *DBWrapper) Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	return db.DB.Query(model, query, params...)
}

// QueryOne acts like Query, but query must return only one row. It
// returns ErrNoRows error when query returns zero rows or
// ErrMultiRows when query returns multiple rows.
func (db *DBWrapper) QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error) {
	return db.DB.QueryOne(model, query, params...)
}
