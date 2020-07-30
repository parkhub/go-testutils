package testutils

import "github.com/go-pg/pg/v9"

type BaseDB interface {
	// RunInTransaction runs a function in a transaction. If function
	// returns an error transaction is rollbacked, otherwise transaction
	// is committed.
	RunInTransaction(fn func(Tx) error) error

	// Query executes a query that returns rows, typically a SELECT.
	// The params are for any placeholders in the query.
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)

	// QueryOne acts like Query, but query must return only one row. It
	// returns ErrNoRows error when query returns zero rows or
	// ErrMultiRows when query returns multiple rows.
	QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)
}

// DB interface includes the pg.DB methods used in transactions API
type DB interface {
	BaseDB

	Model(model ...interface{}) Query
}

// Tx interface includes the pg.Tx methods used in transactions API
type Tx interface {
	// Query runs  an alias for DB.Query
	Query(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)

	// QueryOne is an alias for DB.QueryOne
	QueryOne(model interface{}, query interface{}, params ...interface{}) (pg.Result, error)

	// Select is an alias for DB.Select
	Select(model interface{}) error

	// Insert is an alias for DB.Insert
	Insert(model ...interface{}) error

	// Update is an alias for DB.Update
	Update(model interface{}) error

	// Delete is an alias for DB.Delete
	Delete(model interface{}) error
}
